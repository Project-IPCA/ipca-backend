package student

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceWrapper interface {
	Create(
		userId uuid.UUID,
		KmitlId string,
	) error
}

type Service struct {
	DB *gorm.DB
}

func NewStudentService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
