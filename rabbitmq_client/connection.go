package rabbitmq_client

import (
	"fmt"

	"github.com/Project-IPCA/ipca-backend/config"
	amqp "github.com/rabbitmq/amqp091-go"
)
  
func RabbitMQClient(cfg *config.Config) *amqp.Connection{
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.RABBITMQ.User, cfg.RABBITMQ.Password, cfg.RABBITMQ.Host, cfg.RABBITMQ.Port)
	
	conn, err := amqp.Dial(url)
	if err != nil {
		panic("Failed to connect to RabbitMQ: " + err.Error())
	}

	return conn
}