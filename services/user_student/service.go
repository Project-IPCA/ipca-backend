package userstudent

import (
	"gorm.io/gorm"
)

type ServiceWrapper interface {
	Create(stuId string, firstName string, lastName string) error
}

type Service struct {
	DB *gorm.DB
}

func NewUserStudentService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
