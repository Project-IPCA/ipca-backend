package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Project-IPCA/ipca-backend/models"
)

func (tokenService *Service) CreateRefreshTokenUserStudent(
	user *models.User,
) (t string, err error) {
	claimsRefresh := &JwtCustomRefreshClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * ExpireRefreshCount)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)

	rt, err := refreshToken.SignedString([]byte(tokenService.config.Auth.RefreshSecretUserStudent))
	if err != nil {
		return "", err
	}
	return rt, err
}
