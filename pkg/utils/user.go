package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/Project-IPCA/ipca-backend/repositories"
	"github.com/Project-IPCA/ipca-backend/services/token"
)

func GetUserClaims(c echo.Context, userRepo repositories.UserRepository) (models.User, error) {
	userJwt := c.Get("user").(*jwt.Token)
	claims := userJwt.Claims.(*token.JwtCustomClaims)
	userId := claims.UserID

	existUser := models.User{}
	userRepo.GetUserByUserID(&existUser, userId)
	if claims.CiSession != *existUser.CISession {
		return models.User{}, errors.New("invalid session")
	}
	return existUser, nil
}

func IsRoleSupervisor(user models.User) bool {
	if *user.Role != constants.Role.Supervisor {
		return false
	}
	return true
}
