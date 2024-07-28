package routes

import (
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	s "github.com/Project-IPCA/ipca-backend/server"
	"github.com/Project-IPCA/ipca-backend/server/handlers"
)

func ConfigureRoutes(server *s.Server) {
	testHandler := handlers.NewTestHandler(server)
	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
	server.Echo.Use(middleware.Logger())
	apiGroup := server.Echo.Group("/api")

	apiGroup.GET("/greeting", testHandler.Greeting)
}
