package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/Project-IPCA/ipca-backend/pkg/requests"
	"github.com/Project-IPCA/ipca-backend/pkg/responses"
	"github.com/Project-IPCA/ipca-backend/pkg/utils"
	"github.com/Project-IPCA/ipca-backend/redis_client"
	"github.com/Project-IPCA/ipca-backend/repositories"
	s "github.com/Project-IPCA/ipca-backend/server"
	activitylog "github.com/Project-IPCA/ipca-backend/services/activity_log"
	"github.com/Project-IPCA/ipca-backend/services/token"
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

	if user.UserID == uuid.Nil {
		return responses.ErrorResponse(
			c, http.StatusUnauthorized,
			"username or password is not correct.",
		)
	}

	if *user.Role != constants.Role.Student {
		return responses.ErrorResponse(
			c, http.StatusUnauthorized,
			"username or password is not correct.",
		)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)) != nil {
		return responses.ErrorResponse(
			c,
			http.StatusUnauthorized,
			"username or password is not correct.",
		)
	}

	studentRepository := repositories.NewStudentRepository(authHandler.server.DB)
	classScheduleRepository := repositories.NewClassScheduleRepository(authHandler.server.DB)
	student := models.Student{}
	classSchedule := models.ClassSchedule{}

	studentRepository.GetStudentByStuID(&student, user.UserID)
	classScheduleRepository.GetClassScheduleByGroupID(&classSchedule, *student.GroupID)
	if classSchedule.AllowLogin == false {
		return responses.ErrorResponse(
			c,
			http.StatusForbidden,
			"Login is not allowed by Instructor.",
		)
	}

	activityLogService := activitylog.NewActivityLogService(authHandler.server.DB)

	ip, port, userAgent := utils.GetNetworkRequest(c)

	redis := redis_client.NewRedisAction(authHandler.server.Redis)

	userService := userservice.NewUserService(authHandler.server.DB)
	if user.IsOnline == true {
		userService.UpdateIsOnline(&user, false)

		redisCnl := fmt.Sprintf(
			"%s:%s",
			constants.RedisChannel.UserEvent,
			user.UserID,
		)
		redisMsg := redis.NewMessage("repeat-login", &user.UserID)
		if err := redis.PublishMessage(redisCnl, redisMsg); err != nil {
			return responses.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"Internal Server Error",
			)
		}

		if user.Student != nil {

			redisCnl = fmt.Sprintf(
				"%s:%s",
				constants.RedisChannel.OnlineStudent,
				user.Student.GroupID,
			)
			redisMsg = redis.NewMessage("logout", &user.UserID)
			if err := redis.PublishMessage(redisCnl, redisMsg); err != nil {
				return responses.ErrorResponse(
					c,
					http.StatusInternalServerError,
					"Internal Server Error",
				)
			}

			newLog, err := activityLogService.Create(
				user.Student.GroupID,
				user.Username,
				ip,
				&port,
				&userAgent,
				constants.LogPage.Login,
				constants.LogAction.LoginRepeat,
			)
			if err != nil {
				return responses.ErrorResponse(
					c,
					http.StatusInternalServerError,
					"Can't Insert Log.",
				)
			}

			redisCnl = fmt.Sprintf(
				"%s:%s",
				constants.RedisChannel.Log,
				user.Student.GroupID,
			)
			if err := redis.PublishMessage(redisCnl, newLog); err != nil {
				return responses.ErrorResponse(
					c,
					http.StatusInternalServerError,
					"Internal Server Error",
				)
			}
		}
		return responses.ErrorResponse(
			c,
			http.StatusUnauthorized,
			"Repeat log in. Previous machine logged out. Please try again.",
		)
	}

	userService.UpdateLoginSuccess(&user)

	if user.Student != nil {
		redisCnl := fmt.Sprintf("%s:%s", constants.RedisChannel.OnlineStudent, user.Student.GroupID)
		redisMsg := redis.NewMessage("login", &user.UserID)
		if err := redis.PublishMessage(redisCnl, redisMsg); err != nil {
			return responses.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"Internal Server Error",
			)
		}
	}

	if user.Student != nil {
		newLog, err := activityLogService.Create(
			user.Student.GroupID,
			user.Username,
			ip,
			&port,
			&userAgent,
			constants.LogPage.Login,
			constants.LogAction.Login,
		)
		if err != nil {
			return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Insert Log.")
		}

		redisCnl := fmt.Sprintf(
			"%s:%s",
			constants.RedisChannel.Log,
			user.Student.GroupID,
		)
		if err := redis.PublishMessage(redisCnl, newLog); err != nil {
			return responses.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"Internal Server Error",
			)
		}
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

	response := responses.NewLoginResponse(accessToken, refreshToken, exp)
	return responses.Response(c, http.StatusOK, response)
}

// @Description	Login
// @ID				auth-login-super
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			params	body		requests.LoginRequest	true	"User's credentials"
// @Success		200		{object}	responses.LoginResponse
// @Failure		401		{object}	responses.Error
// @Router			/api/auth/login/super [post]
func (authHandler *AuthHandler) LoginSuper(c echo.Context) error {
	loginReq := new(requests.LoginRequest)

	if err := c.Bind(loginReq); err != nil {
		return err
	}

	user := models.User{}
	userRepository := repositories.NewUserRepository(authHandler.server.DB)
	userRepository.GetUserByUsername(&user, loginReq.Username)

	if user.UserID == uuid.Nil {
		return responses.ErrorResponse(
			c, http.StatusNotFound,
			"User Not Found.",
		)
	}

	if *user.Role == constants.Role.Student {
		return responses.ErrorResponse(
			c, http.StatusUnauthorized,
			"username or password is not correct.",
		)
	}

	if !user.IsActive {
		return responses.ErrorResponse(
			c, http.StatusForbidden,
			"This Admin Has Been Deleted.",
		)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)) != nil {
		return responses.ErrorResponse(
			c,
			http.StatusUnauthorized,
			"username or password is not correct.",
		)
	}
	redis := redis_client.NewRedisAction(authHandler.server.Redis)

	userService := userservice.NewUserService(authHandler.server.DB)
	if user.IsOnline == true {
		userService.UpdateIsOnline(&user, false)

		redisCnl := fmt.Sprintf(
			"%s:%s",
			constants.RedisChannel.UserEvent,
			user.UserID,
		)
		redisMsg := redis.NewMessage("repeat-login", &user.UserID)
		if err := redis.PublishMessage(redisCnl, redisMsg); err != nil {
			return responses.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"Internal Server Error",
			)
		}

		return responses.ErrorResponse(
			c,
			http.StatusUnauthorized,
			"Repeat log in. Previous machine logged out. Please try again.",
		)
	}

	userService.UpdateLoginSuccess(&user)

	tokenService := tokenservice.NewTokenService(authHandler.server.Config)
	accessToken, exp, err := tokenService.CreateAccessToken(&user)
	if err != nil {
		return err
	}
	refreshToken, err := tokenService.CreateRefreshToken(&user)
	if err != nil {
		return err
	}

	response := responses.NewLoginResponse(accessToken, refreshToken, exp)
	return responses.Response(c, http.StatusOK, response)
}

// @Description	Logout
// @ID				auth-logout
// @Tags			Auth
// @Accept		json
// @Produce		json
// @Success		200		{object}	responses.Data
// @Failure		404		{object}	responses.Error
// @Failure		500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/auth/logout [post]
func (authHandler *AuthHandler) Logout(c echo.Context) error {
	userJwt := c.Get("user").(*jwt.Token)
	claims := userJwt.Claims.(*token.JwtCustomClaims)
	userId := claims.UserID

	userService := userservice.NewUserService(authHandler.server.DB)
	existsUser := models.User{}

	userRepository := repositories.NewUserRepository(authHandler.server.DB)
	userRepository.GetUserByUserID(&existsUser, userId)

	if existsUser.UserID != userId {
		return responses.ErrorResponse(c, http.StatusNotFound, "User ID not found in session.")
	}

	userService.UpdateIsOnline(&existsUser, false)

	if *existsUser.Role == constants.Role.Student {
		redis := redis_client.NewRedisAction(authHandler.server.Redis)
		redisCnl := fmt.Sprintf(
			"%s:%s",
			constants.RedisChannel.OnlineStudent,
			existsUser.Student.GroupID,
		)
		redisMsg := redis.NewMessage("logout", &existsUser.UserID)
		if err := redis.PublishMessage(redisCnl, redisMsg); err != nil {
			return responses.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"Internal Server Error",
			)
		}
		activityLogService := activitylog.NewActivityLogService(authHandler.server.DB)
		ip, port, userAgent := utils.GetNetworkRequest(c)
		newLog, err := activityLogService.Create(
			existsUser.Student.GroupID,
			existsUser.Username,
			ip,
			&port,
			&userAgent,
			constants.LogPage.Login,
			constants.LogAction.Logout,
		)
		
		if err != nil {
			return responses.ErrorResponse(c, http.StatusInternalServerError, "Can't Insert Log.")
		}

		redisCnl = fmt.Sprintf(
			"%s:%s",
			constants.RedisChannel.Log,
			existsUser.Student.GroupID,
		)
		if err := redis.PublishMessage(redisCnl, newLog); err != nil {
			return responses.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"Internal Server Error",
			)
		}
		
	}

	return responses.MessageResponse(c, http.StatusOK, "Logout successful")
}

// @Description	Refresh Token
// @ID				auth-refresh-token
// @Tags			Auth
// @Accept		json
// @Produce		json
// @Success		200		{object}	responses.LoginResponse
// @Failure		403		{object}	responses.Error
// @Failure		500		{object}	responses.Error
// @Security BearerAuth
// @Router			/api/auth/refresh_token [post]
func (authHandler *AuthHandler) RefreshToken(c echo.Context) error {
	userJwt := c.Get("user").(*jwt.Token)
	claims := userJwt.Claims.(*token.JwtCustomRefreshClaims)
	userId := claims.UserID

	var existsUser models.User
	userRepository := repositories.NewUserRepository(authHandler.server.DB)
	userRepository.GetUserByUserID(&existsUser, userId)

	if claims.CiSession != *existsUser.CISession {
		return responses.ErrorResponse(c, http.StatusForbidden, "Invalid Session.")
	}

	if !existsUser.IsActive {
		return responses.ErrorResponse(
			c, http.StatusForbidden,
			"This Admin Has Been Deleted.",
		)
	}

	tokenService := tokenservice.NewTokenService(authHandler.server.Config)
	accessToken, exp, err := tokenService.CreateAccessToken(&existsUser)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	refreshToken, err := tokenService.CreateRefreshToken(&existsUser)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	response := responses.NewLoginResponse(accessToken, refreshToken, exp)
	return responses.Response(c, http.StatusOK, response)
}
