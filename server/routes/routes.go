package routes

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/Project-IPCA/ipca-backend/middlewares"
	s "github.com/Project-IPCA/ipca-backend/server"
	"github.com/Project-IPCA/ipca-backend/server/handlers"
)

func ConfigureRoutes(server *s.Server) {
	testHandler := handlers.NewTestHandler(server)
	supervisorHandler := handlers.NewSupervisorHandler(server)
	studentHandler := handlers.NewStudentHandle(server)
	authHandler := handlers.NewAuthHandler(server)
	commonHandler := handlers.NewCommonHandler(server)
	initHandler := handlers.NewInitHandler(server)

	server.Echo.IPExtractor = echo.ExtractIPFromXFFHeader()
	server.Echo.Static("/static", "bucket")

	server.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
	server.Echo.Use(middleware.Logger())
	server.Echo.Use(middleware.CORS())

	authMiddleware := middlewares.NewAuthMiddleware(server)
	jwtConfig := authMiddleware.GetJwtConfig()
	refreshTokenConfig := authMiddleware.GetRefreshTokenConfig()

	apiGroup := server.Echo.Group("/api")

	// Init
	initGroup := apiGroup.Group("/init")
	initGroup.POST("/department", initHandler.InitDepartment)
	initGroup.POST("/supervisor", initHandler.InitSupervisor)
	initGroup.POST("/labclassinfo", initHandler.InitClassInfo)
	initGroup.POST("/ta", initHandler.InitTA)

	// Supervisor
	supervisorGroup := apiGroup.Group("/supervisor")
	supervisorAuthGroup := supervisorGroup
	supervisorAuthGroup.Use(echojwt.WithConfig(jwtConfig))
	supervisorAuthGroup.POST("/students", supervisorHandler.AddStudents)
	supervisorAuthGroup.POST("/group", supervisorHandler.CreateGroup)
	supervisorAuthGroup.DELETE("/group/:group_id", supervisorHandler.DeleteGroup)
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
	supervisorAuthGroup.POST("/save_exercise_testcase", supervisorHandler.SaveExerciseTestcase)
	supervisorAuthGroup.POST(
		"/update_group_assigned_chapter_item",
		supervisorHandler.UpdateGroupAssignedChapterItem,
	)
	supervisorAuthGroup.POST(
		"/update_all_group_assigned_chapter_item",
		supervisorHandler.UpdateAllGroupAssignedChapterItem,
	)
	supervisorAuthGroup.GET("/get_lab_chapter_info", supervisorHandler.GetLabChapterInfo)
	supervisorAuthGroup.GET("/get_student_group_list", supervisorHandler.GetStudentGroupList)
	supervisorAuthGroup.POST("/set_chapter_permission", supervisorHandler.SetChapterPemission)
	supervisorAuthGroup.POST("/set_allow_group_login", supervisorHandler.SetAllowGroupLogin)
	supervisorAuthGroup.POST(
		"/set_allow_group_upload_picture",
		supervisorHandler.SetAllowGroupUploadPicture,
	)
	supervisorAuthGroup.DELETE("/exercise/:exercise_id", supervisorHandler.DeleteExercise)
	supervisorAuthGroup.PUT(
		"/student_can_submit/:student_id",
		supervisorHandler.UpdateStudentCanSubmit,
	)
	supervisorAuthGroup.GET("/student_info/:student_id", supervisorHandler.GetStudentInfo)
	supervisorAuthGroup.PUT(
		"/logout_all_student/:group_id",
		supervisorHandler.LogoutAllStudentInGroup,
	)
	supervisorAuthGroup.DELETE("/student/:stu_id", supervisorHandler.DeleteStudent)
	supervisorAuthGroup.GET("/get_exercise_data/:exercise_id", supervisorHandler.GetExerciseData)
	supervisorAuthGroup.PUT(
		"/cancle_student_submission/:submission_id",
		supervisorHandler.CancleStduentSubmission,
	)
	supervisorAuthGroup.GET(
		"/student_chapter_list",
		supervisorHandler.GetStudentChapterList,
	)
	supervisorAuthGroup.GET("/assigned_student_exercise", supervisorHandler.GetAssginStudentExercise)
	supervisorAuthGroup.PUT("/exercise", supervisorHandler.UpdateExercise)
	supervisorAuthGroup.POST("/admin", supervisorHandler.CreateAdmin)
	supervisorAuthGroup.POST("/department", supervisorHandler.CreateDepartment)
	supervisorAuthGroup.POST("/set_role_permission", supervisorHandler.SetRolePermission)
	supervisorAuthGroup.GET("/role_permission", supervisorHandler.GetRolePermission)

	// Student
	studentGroup := apiGroup.Group("/student")
	studentGroup.Use(echojwt.WithConfig(jwtConfig))
	studentGroup.POST("/exercise_submit", studentHandler.ExerciseSubmit)
	studentGroup.GET("/all_chapter", studentHandler.GetALLChapter)
	studentGroup.GET("/chapter_list", studentHandler.GetChapterList)
	studentGroup.GET("/assigned_exercise", studentHandler.GetStudentAssignedExercise)

	// Auth
	authGroup := apiGroup.Group("/auth")
	authGroup.POST("/login", authHandler.Login)
	authAuthGroup := authGroup
	authAuthGroup.Use(echojwt.WithConfig(jwtConfig))
	authAuthGroup.POST("/logout", authHandler.Logout)
	refreshAuthGroup := authGroup
	refreshAuthGroup.Use(echojwt.WithConfig(refreshTokenConfig))
	refreshAuthGroup.POST("/refresh_token", authHandler.RefreshToken)

	// Common
	commonGroup := apiGroup.Group("/common")
	commonAuthGroup := commonGroup
	commonAuthGroup.Use(echojwt.WithConfig(jwtConfig))
	commonAuthGroup.GET("/user_info", commonHandler.GetUserInfo)
	commonAuthGroup.PUT("/user_info", commonHandler.UpdateUserInfo)
	commonAuthGroup.POST("/get_keyword_list", commonHandler.GetKeywordList)
	commonAuthGroup.POST("/keyword_check", commonHandler.KeywordCheck)
	commonAuthGroup.GET("/student_submission", commonHandler.GetStudentSubmission)
	commonAuthGroup.POST("/user_profile", commonHandler.UploadUserProfile)
	commonAuthGroup.GET("/departments", commonHandler.GetDepartments)
	commonAuthGroup.GET("/staffs", commonHandler.GetStaffs)

	// Test
	apiGroup.GET("/greeting", testHandler.Greeting)
	apiGroup.GET("/test_redis", testHandler.TestRedis)
	apiGroup.GET("/test_rabbit", testHandler.TestRabbitMQ)
}
