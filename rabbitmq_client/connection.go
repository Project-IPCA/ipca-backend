package rabbitmq_client

import (
	"fmt"

	"github.com/Project-IPCA/ipca-backend/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConnection struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func RabbitMQClient(cfg *config.Config) *RabbitMQConnection {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)

	conn, err := amqp.Dial(url)
	if err != nil {
		panic("Failed to connect to RabbitMQ: " + err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		panic("Fail to open channel: " + err.Error())
	}

	return &RabbitMQConnection{
		Conn: conn,
		Ch:   ch,
	}
}
