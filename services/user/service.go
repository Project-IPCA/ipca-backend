package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceWrapper interface {
	CreateQuick(
		username string,
		password string,
		firstName string,
		lastName string,
		role string,
	) (userID uuid.UUID, err error)
}

type Service struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
