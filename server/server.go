package server

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/config"
	"github.com/Project-IPCA/ipca-backend/db"
)

type Server struct {
	Echo   *echo.Echo
	DB     *gorm.DB
	Config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Echo:   echo.New(),
		DB:     db.Init(cfg),
		Config: cfg,
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}
