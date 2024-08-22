package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type LabExerciseRepositoryQ interface{}

type LabExerciseRepository struct {
	DB *gorm.DB
}

func NewLabExerciseRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{DB: db}
}

func (labExerciesRepository *LabExerciseRepository) GetLabExerciseOrderByChapterLevel() {
	labExerciesRepository.DB.Preload("Chapter")
}
