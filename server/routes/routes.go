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
	commonHandler := handlers.NewCommonHandler(server)
	initHandler := handlers.NewInitHandler(server)
  

	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
	server.Echo.Use(middleware.Logger())
  server.Echo.Use(middleware.CORS())

	authMiddleware := middlewares.NewAuthMiddleware(server)
	jwtConfig := authMiddleware.GetJwtConfig()

	apiGroup := server.Echo.Group("/api")

	// Init
	initGroup := apiGroup.Group("/init")
	initGroup.POST("/department", initHandler.InitDepartment)
	initGroup.POST("/supervisor", initHandler.InitSupervisor)

	// Supervisor
	supervisorGroup := apiGroup.Group("/supervisor")
	supervisorAuthGroup := supervisorGroup
	supervisorAuthGroup.Use(echojwt.WithConfig(jwtConfig))
	supervisorAuthGroup.POST("/students", supervisorHandler.AddStudents)
	supervisorAuthGroup.POST("/group", supervisorHandler.CreateGroup)
	supervisorAuthGroup.POST("/exercise", supervisorHandler.CreateExercise)
	supervisorAuthGroup.GET("/available_groups", supervisorHandler.GetAllAvailableGroups)
	supervisorAuthGroup.GET("/my_groups", supervisorHandler.GetMyGroups)
	supervisorAuthGroup.GET("/group/:group_id", supervisorHandler.GetGroupInfoByGroupID)
	supervisorAuthGroup.PUT(
		"/reset_student_password/:stu_id",
		supervisorHandler.ResetStudentPassword,
	)
	supervisorAuthGroup.GET("/my_group_info/:group_id", supervisorHandler.GetMyGroupInfo)
	supervisorAuthGroup.PUT("/my_group_info/:group_id", supervisorHandler.UpdateMyGroupInfo)

	// Auth
	authGroup := apiGroup.Group("/auth")
	authGroup.POST("/login", authHandler.Login)
	authAuthGroup := authGroup
	authAuthGroup.Use(echojwt.WithConfig(jwtConfig))
	authAuthGroup.POST("/logout", authHandler.Logout)

	// Common
	commonGroup := apiGroup.Group("/common")
	commonAuthGroup := commonGroup
	commonAuthGroup.Use(echojwt.WithConfig(jwtConfig))
	commonAuthGroup.GET("/user_info", commonHandler.GetUserInfo)
	commonAuthGroup.PUT("/user_info", commonHandler.UpdateUserInfo)
	commonAuthGroup.POST("/get_keyword_list", commonHandler.GetKeywordList)
	commonAuthGroup.POST("/keyword_check", commonHandler.KeywordCheck)

	// Test
	apiGroup.GET("/greeting", testHandler.Greeting)
	apiGroup.GET("/test_redis", testHandler.TestRedis)
	apiGroup.GET("/test_rabbit", testHandler.TestRabbitMQ)
}
