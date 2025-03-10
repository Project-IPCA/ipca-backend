package exercisetestcase

import (
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/pkg/requests"
)

type ServiceWrapper interface {
	Create(request *requests.CreateGroupRequest) error
}

type Service struct {
	DB *gorm.DB
}

func NewExerciseTestcaseService(db *gorm.DB) *Service {
	return &Service{DB: db}
}