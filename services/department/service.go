package department

import (
	"gorm.io/gorm"
)

type ServiceWrapper interface {
	Create(name string) error
}

type Service struct {
	DB *gorm.DB
}

func NewDepartmetService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
