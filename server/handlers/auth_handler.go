package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/Project-IPCA/ipca-backend/pkg/requests"
	"github.com/Project-IPCA/ipca-backend/pkg/responses"
	"github.com/Project-IPCA/ipca-backend/repositories"
	s "github.com/Project-IPCA/ipca-backend/server"
	tokenservice "github.com/Project-IPCA/ipca-backend/services/token"
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
func (authHandler AuthHandler) Login(c echo.Context) error {
	loginReq := new(requests.LoginRequest)

	if err := c.Bind(loginReq); err != nil {
		return err
	}

	user := models.User{}
	userRepository := repositories.NewUserRepository(authHandler.server.DB)
	userRepository.GetUserByUsername(&user, loginReq.Username)

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)) != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Invalid Credentials")
	}

	tokenService := tokenservice.NewTokenService(authHandler.server.Config)
	if user.Role == constants.Role.Student {
		accessToken, exp, err := tokenService.CreateUserStudentAccessToken(&user)
		if err != nil {
			return err
		}
		refreshToken, err := tokenService.CreateRefreshTokenUserStudent(&user)
		if err != nil {
			return err
		}
		response := responses.NewLoginResponse(accessToken, refreshToken, exp)
		return responses.Response(c, http.StatusOK, response)
	}
	return nil
}
