package rabbitmq_client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Project-IPCA/ipca-backend/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type IRabbitMQAction interface {
	SendQueue(message string)
}

type RabbitMQAction struct {
	RabitMQ *RabbitMQConnection
	cfg     *config.Config
}

func NewRabbitMQAction(rabbitmq *RabbitMQConnection, config *config.Config) *RabbitMQAction {
	return &RabbitMQAction{
		RabitMQ: rabbitmq,
		cfg:     config,
	}
}

func (rabbitMQAction *RabbitMQAction) SendQueue(message interface{}) error {
	var ch *amqp.Channel
	maxRetries := 5
	for i := 1; i <= maxRetries; i++ {
		tempCh, err := rabbitMQAction.GetChannel()
		if err != nil {
			reconnectErr := rabbitMQAction.Reconnect()
			if reconnectErr != nil {
				if i == maxRetries {
					return fmt.Errorf("fail to reconnect: %v", reconnectErr)
				}
				time.Sleep(time.Duration(i+1) * time.Second)
				continue
			}

			tempCh, err = rabbitMQAction.GetChannel()
			if err != nil {
				if i == maxRetries {
					return fmt.Errorf("fail to get channel after reconnect: %v", reconnectErr)
				}
				time.Sleep(time.Duration(i+1) * time.Second)
				continue
			}
		}
		ch = tempCh
		break
	}

	q, err := ch.QueueDeclare(
		rabbitMQAction.cfg.RabbitMQ.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("fail declare a queue: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message to json: %v", err)
	}

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to send queue to rabbitmq: %v", err)
	}

	return nil
}

func (rmq *RabbitMQAction) GetChannel() (*amqp.Channel, error) {
	rmq.RabitMQ.Mu.Lock()
	defer rmq.RabitMQ.Mu.Unlock()

	if rmq.RabitMQ.Ch.IsClosed() {
		ch, err := rmq.RabitMQ.Conn.Channel()
		if err != nil {
			return nil, fmt.Errorf("failed to create new channel: %v", err)
		}
		rmq.RabitMQ.Ch = ch
	}

	return rmq.RabitMQ.Ch, nil
}

func (rmq *RabbitMQAction) Reconnect() error {
	rmq.RabitMQ.Mu.Lock()
	defer rmq.RabitMQ.Mu.Unlock()

	if rmq.RabitMQ.Ch != nil {
		rmq.RabitMQ.Ch.Close()
	}
	if rmq.RabitMQ.Conn != nil {
		rmq.RabitMQ.Conn.Close()
	}

	conn, err := amqp.Dial(rmq.RabitMQ.Url)
	if err != nil {
		return fmt.Errorf("failed to reconnect: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to create channel: %v", err)
	}

	rmq.RabitMQ.Conn = conn
	rmq.RabitMQ.Ch = ch

	return nil
}
