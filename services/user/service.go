package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type ServiceWrapper interface {
	CreateQuick(
		username string,
		password string,
		firstName string,
		lastName string,
		role string,
	) (userID uuid.UUID, err error)
	UpdateLoginSuccess(user *models.User) error
	UpdateIsOnline(user *models.User, isOnline bool) error
}

type Service struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
