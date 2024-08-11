package rabbitmq_client

import (
	"os"
	"context"
	"encoding/json"
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

func NewRabbitMQAction(rabbitmq *amqp.Connection) *RabbitMQAction {
	return &RabbitMQAction{
		RabbitMQ: rabbitmq,
	}
}

func (rabbitMQAction *RabbitMQAction) SendQueue(message interface{}) error {
	ch, err := rabbitMQAction.RabbitMQ.Channel()
	if err != nil {
		panic("fail to open channel: " + err.Error())
	}

	q, err := ch.QueueDeclare(
		os.Getenv("RABBITMQ_QUEUENAME"),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("fail declare a queue: " + err.Error())
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(message)
	if err != nil {
		panic("failed to marshal message to JSON: " + err.Error())
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
		panic("failed to send queue to RabbitMQ: " + err.Error())
	}

	return nil
}