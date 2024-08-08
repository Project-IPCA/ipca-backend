package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Project-IPCA/ipca-backend/models"
)

func (tokenService *Service) CreateUserStudentAccessToken(
	user *models.User,
) (t string, expired int64, err error) {
	exp := time.Now().Add(time.Hour * ExpireCount)
	claims := &JwtCustomClaims{
		user.ID,
		user.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	expired = exp.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err = token.SignedString([]byte(tokenService.config.Auth.AccessSecretUserStudent))
	if err != nil {
		return
	}

	return
}
