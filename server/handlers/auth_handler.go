package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/Project-IPCA/ipca-backend/pkg/requests"
	"github.com/Project-IPCA/ipca-backend/pkg/responses"
	"github.com/Project-IPCA/ipca-backend/redis_client"
	"github.com/Project-IPCA/ipca-backend/repositories"
	s "github.com/Project-IPCA/ipca-backend/server"
	tokenservice "github.com/Project-IPCA/ipca-backend/services/token"
	userservice "github.com/Project-IPCA/ipca-backend/services/user"
)

type AuthHandler struct {
	server *s.Server
}

func NewAuthHandler(server *s.Server) *AuthHandler {
	return &AuthHandler{server: server}
}

// @Description	Login
// @ID				auth-login
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			params	body		requests.LoginRequest	true	"User's credentials"
// @Success		200		{object}	responses.LoginResponse
// @Failure		401		{object}	responses.Error
// @Router			/api/auth/login [post]
func (authHandler *AuthHandler) Login(c echo.Context) error {
	loginReq := new(requests.LoginRequest)

	if err := c.Bind(loginReq); err != nil {
		return err
	}

	user := models.User{}
	userRepository := repositories.NewUserRepository(authHandler.server.DB)
	userRepository.GetUserByUsername(&user, loginReq.Username)

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)) != nil {
		return responses.ErrorResponse(
			c,
			http.StatusUnauthorized,
			"username or password is not correct.",
		)
	}

	tokenService := tokenservice.NewTokenService(authHandler.server.Config)
	accessToken, exp, err := tokenService.CreateAccessToken(&user)
	if err != nil {
		return err
	}
	refreshToken, err := tokenService.CreateRefreshToken(&user)
	if err != nil {
		return err
	}

	studentRepository := repositories.NewStudentRepository(authHandler.server.DB)
	classScheduleRepository := repositories.NewClassScheduleRepository(authHandler.server.DB)
	student := models.Student{}
	classSchedule := models.ClassSchedule{}

	if user.Role == &constants.Role.Student {
		studentRepository.GetStudentByStuID(&student, user.UserID)
		classScheduleRepository.GetClassScheduleByGroupID(&classSchedule, *student.GroupID)
		if classSchedule.AllowLogin == false {
			return responses.ErrorResponse(
				c,
				http.StatusUnauthorized,
				"Login is not allowed by Instructor.",
			)
		}
	}

	userService := userservice.NewUserService(authHandler.server.DB)
	if user.IsOnline == true {
		userService.UpdateIsOnline(&user, false)
		return responses.ErrorResponse(
			c,
			http.StatusUnauthorized,
			"Repeat log in. Previous machine logged out. Please try again.",
		)
	}

	userService.UpdateLoginSuccess(&user)

	redis := redis_client.NewRedisAction(authHandler.server.Redis)
	redisCnl := fmt.Sprintf("online-students:%s", user.UserID)
	redisMsg := redis.NewMessage("login", user.UserID)
	if err := redis.PublishMessage(redisCnl, redisMsg); err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
	}

	response := responses.NewLoginResponse(accessToken, refreshToken, exp)
	return responses.Response(c, http.StatusOK, response)
}
