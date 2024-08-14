package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/Project-IPCA/ipca-backend/middlewares"
	s "github.com/Project-IPCA/ipca-backend/server"
	"github.com/Project-IPCA/ipca-backend/server/handlers"
)

func ConfigureRoutes(server *s.Server) {
	testHandler := handlers.NewTestHandler(server)
	supervisorHandler := handlers.NewSupervisorHandler(server)
	authHandler := handlers.NewAuthHandler(server)

	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
	server.Echo.Use(middleware.Logger())

	authMiddleware := middlewares.NewAuthMiddleware(server)
	jwtConfig := authMiddleware.GetJwtConfig()

	apiGroup := server.Echo.Group("/api")

	// Supervisor
	supervisorGroup := apiGroup.Group("/supervisor")
	supervisorAuthGroup := supervisorGroup
	supervisorAuthGroup.Use(echojwt.WithConfig(jwtConfig))
	supervisorAuthGroup.POST("/add_students", supervisorHandler.AddStudents)

	// Auth
	authGroup := apiGroup.Group("/auth")
	authGroup.POST("/login", authHandler.Login)
	authAuthGroup := authGroup
	authAuthGroup.Use(echojwt.WithConfig(jwtConfig))
	authAuthGroup.POST("/logout", authHandler.Logout)

	// Test
	apiGroup.GET("/greeting", testHandler.Greeting)
	apiGroup.GET("/test_redis", testHandler.TestRedis)
	apiGroup.GET("/test_rabbit", testHandler.TestRabbitMQ)
}
