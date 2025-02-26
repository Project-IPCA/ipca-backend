package server

import (
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/config"
	"github.com/Project-IPCA/ipca-backend/db"
	minioclient "github.com/Project-IPCA/ipca-backend/minio_client"
	"github.com/Project-IPCA/ipca-backend/rabbitmq_client"
	"github.com/Project-IPCA/ipca-backend/redis_client"
)

type Server struct {
	Echo    *echo.Echo
	DB      *gorm.DB
	Config  *config.Config
	Redis   *redis.Client
	RabitMQ *rabbitmq_client.RabbitMQConnection
	Minio   *minio.Client
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Echo:    echo.New(),
		DB:      db.Init(cfg),
		Config:  cfg,
		Redis:   redis_client.RedisClient(cfg),
		RabitMQ: rabbitmq_client.RabbitMQClient(cfg),
		Minio:   minioclient.Init(cfg),
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}

func (server *Server) CloseConnection() {
	if server.RabitMQ != nil {
		if server.RabitMQ.Ch != nil {
			server.RabitMQ.Ch.Close()
		}
		if server.RabitMQ.Conn != nil {
			server.RabitMQ.Conn.Close()
		}
	}
	if server.DB != nil {
		database, _ := server.DB.DB()
		database.Close()
	}
	if server.Redis != nil {
		server.Redis.Close()
	}
}
