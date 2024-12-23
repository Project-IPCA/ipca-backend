package ta

import (
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewTaService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
