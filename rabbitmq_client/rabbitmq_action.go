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
	RabbitMQ *amqp.Connection
	cfg      *config.Config
}

func NewRabbitMQAction(rabbitmq *amqp.Connection,config *config.Config) *RabbitMQAction {
	return &RabbitMQAction{
		RabbitMQ: rabbitmq,
		cfg: config,
	}
}

func (rabbitMQAction *RabbitMQAction) SendQueue(message interface{}) error {
	ch, err := rabbitMQAction.RabbitMQ.Channel()
	if err != nil {
		return fmt.Errorf("fail to open channel: %v " , err)
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
		return fmt.Errorf("fail declare a queue: %v" ,err)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message to json: %v" ,err)
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