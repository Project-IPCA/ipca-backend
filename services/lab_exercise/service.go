package labexercise

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/pkg/requests"
)

type ServiceWrapper interface {
	CreatePythonExercise(
		request *requests.CreatePythonExerciseRequest,
		supervisorId *uuid.UUID,
		supervisorName string,
	) error
}

type Service struct {
	DB *gorm.DB
}

func NewLabExerciseService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
