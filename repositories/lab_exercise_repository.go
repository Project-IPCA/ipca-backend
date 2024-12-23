package repositories

import (
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
)

type LabExerciseRepositoryQ interface {
	GetLabExerciseOrderByChapterLevel(
		labExercises *[]models.LabExercise,
	)
}

type LabExerciseRepository struct {
	DB *gorm.DB
}

func NewLabExerciseRepository(db *gorm.DB) *LabExerciseRepository {
	return &LabExerciseRepository{DB: db}
}

func (labExerciesRepository *LabExerciseRepository) GetLabExerciseOrderByChapterLevel(
	labExercises *[]models.LabExercise,
) {
	labExerciesRepository.DB.
		Preload("Chapter").
		Joins("LEFT JOIN lab_class_infos ON lab_class_infos.chapter_id = lab_exercises.chapter_id").
		Order("lab_class_infos.chapter_index ASC, lab_exercises.level ASC").
		Find(labExercises)
}

func (labExerciesRepository *LabExerciseRepository) GetLabExerciseByID(exerciseId uuid.UUID, labExercise *models.LabExercise) error {
	err := labExerciesRepository.DB.Preload("Chapter").Preload("TestcaseList").Where("exercise_id = ?", exerciseId).First(labExercise)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (labExerciesRepository *LabExerciseRepository) GetLabExerciseByChapterID(labExercise *[]models.LabExercise, chapterId uuid.UUID) {
	labExerciesRepository.DB.Where("chapter_id = ?", chapterId).Find(labExercise)
}
