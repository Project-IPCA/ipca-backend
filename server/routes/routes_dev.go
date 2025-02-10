//go:build dev

package routes

import (
	"github.com/Project-IPCA/ipca-backend/server"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// ConfigureSwaggerRoutes enables Swagger for development mode.
func ConfigureSwaggerRoutes(server *server.Server) {
	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
}
