package executive

import (
	"gorm.io/gorm"
)

type ServiceWrapper interface{}

type Service struct {
	DB *gorm.DB
}

func NewExecutiveService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
