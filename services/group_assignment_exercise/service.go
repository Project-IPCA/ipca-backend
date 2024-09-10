package groupassignmentexercise

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceWrapper interface {
	Create(
		groupId uuid.UUID,
		exerciseId uuid.UUID,
		selected bool,
	) error
}

type Service struct {
	DB *gorm.DB
}

func NewGroupAssignmentExerciseService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
