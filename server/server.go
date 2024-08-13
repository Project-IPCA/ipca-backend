package server

import (
	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/config"
	"github.com/Project-IPCA/ipca-backend/db"
	"github.com/Project-IPCA/ipca-backend/redis_client"
)

type Server struct {
	Echo    *echo.Echo
	DB      *gorm.DB
	Config  *config.Config
	Redis   *redis.Client
	RabitMQ *amqp.Connection
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Echo:   echo.New(),
		DB:     db.Init(cfg),
		Config: cfg,
		Redis:  redis_client.RedisClient(cfg),
		// RabitMQ: rabbitmq_client.RabbitMQClient(cfg),
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}

// test github
