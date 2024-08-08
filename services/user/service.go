package user

import (
	"gorm.io/gorm"
)

type ServiceWrapper interface {
	Create(stuId string, role string) error
}

type Service struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
