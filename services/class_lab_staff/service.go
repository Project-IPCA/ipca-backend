package classlabstaff

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceWrapper interface {
	Create(classId uuid.UUID, staffId uuid.UUID) error
}

type Service struct {
	DB *gorm.DB
}

func NewClassLabStaffService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
