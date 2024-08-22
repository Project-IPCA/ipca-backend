package supervisor

import (
	"gorm.io/gorm"
)

type ServiceWrapper interface{}

type Service struct {
	DB *gorm.DB
}

func NewSupervisorService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
