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
	return &ExerciseSubmissionRepository{DB: db}
}

func (exerciseSubmissionRepo *ExerciseSubmissionRepository) GetStudentSubmission(stuId uuid.UUID, exerciseId uuid.UUID, submission *[]models.ExerciseSubmission) {
	exerciseSubmissionRepo.DB.Where("stu_id = ? AND exercise_id = ?", stuId, exerciseId).Order("time_submit ASC").Find(submission)
}

func (exerciseSubmissionRepo *ExerciseSubmissionRepository) GetSubmissionByExerciseID(exerciseId uuid.UUID, submission *[]models.ExerciseSubmission) {
	exerciseSubmissionRepo.DB.Where("exercise_id = ?", exerciseId).Find(submission)
}

func (exerciseSubmissionRepo *ExerciseSubmissionRepository) GetSubmissionByStudentID(stuId uuid.UUID, submission *[]models.ExerciseSubmission) {
	exerciseSubmissionRepo.DB.Where("stu_id = ?", stuId).Find(submission)
}

func (exerciseSubmissionRepo *ExerciseSubmissionRepository) GetSubmissionByID(submissionId uuid.UUID, submission *models.ExerciseSubmission) {
	exerciseSubmissionRepo.DB.Where("submission_id = ?", submissionId).Preload("LabExercise").Preload("Student").Find(submission)
}
