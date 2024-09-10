package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/config"
	"github.com/Project-IPCA/ipca-backend/models"
)

const (
	ExpireCount        = 2
	ExpireRefreshCount = 168
)

type JwtCustomClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type ServiceWrapper interface {
	CreateUserStudentAccessToken(
		user *models.User,
	) (t string, expired int64, err error)
}

type Service struct {
	config *config.Config
}

func NewTokenService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}
