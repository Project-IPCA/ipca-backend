package repositories

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExerciseSubmissionRepository struct {
	DB *gorm.DB
}

func NewExerciseSubmissionRepository(db *gorm.DB) *ExerciseSubmissionRepository {
	return &ExerciseSubmissionRepository{DB:db}
}

func (exerciseSubmissionRepo *ExerciseSubmissionRepository) GetStudentSubmission(stuId uuid.UUID,exerciseId uuid.UUID,submission *[]models.ExerciseSubmission){
	exerciseSubmissionRepo.DB.Where("stu_id = ? AND exercise_id = ?",stuId,exerciseId).Order("time_submit ASC").Find(submission)
}