package server

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/config"
	"github.com/Project-IPCA/ipca-backend/db"
	"github.com/Project-IPCA/ipca-backend/redis_client"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	Echo   *echo.Echo
	DB     *gorm.DB
	Config *config.Config
	Redis 	*redis.Client
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Echo:   echo.New(),
		DB:     db.Init(cfg),
		Config: cfg,
		Redis: redis_client.RedisClient(cfg),
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}
