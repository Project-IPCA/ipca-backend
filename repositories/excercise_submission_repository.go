package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type ExerciseSubmissionRepository struct {
	DB *gorm.DB
}

func NewExerciseSubmissionRepository(db *gorm.DB) *ExerciseSubmissionRepository {
	return &ExerciseSubmissionRepository{DB: db}
}

func (exerciseSubmissionRepo *ExerciseSubmissionRepository) GetStudentSubmission(
	stuId uuid.UUID,
	exerciseId uuid.UUID,
	submission *[]models.ExerciseSubmission,
) {
	exerciseSubmissionRepo.DB.Where("stu_id = ? AND exercise_id = ?", stuId, exerciseId).
		Order("time_submit ASC").
		Find(submission)
}

func (exerciseSubmissionRepo *ExerciseSubmissionRepository) GetSubmissionByExerciseID(
	exerciseId uuid.UUID,
	submission *[]models.ExerciseSubmission,
) {
	exerciseSubmissionRepo.DB.Where("exercise_id = ?", exerciseId).Find(submission)
}

func (exerciseSubmissionRepo *ExerciseSubmissionRepository) GetSubmissionByStudentID(
	stuId uuid.UUID,
	submission *[]models.ExerciseSubmission,
) {
	exerciseSubmissionRepo.DB.Where("stu_id = ?", stuId).Find(submission)
}

func (exerciseSubmissionRepo *ExerciseSubmissionRepository) GetSubmissionByID(
	submissionId uuid.UUID,
	submission *models.ExerciseSubmission,
) error {
	err := exerciseSubmissionRepo.DB.Where("submission_id = ?", submissionId).
		Preload("LabExercise").
		Preload("Student").
		First(submission)
	if err != nil {
		return err.Error
	}
	return nil
}

func (exerciseSubmissionRepo *ExerciseSubmissionRepository) GetTotalSubmissions(
	groupId *uuid.UUID,
) int64 {
	var total int64
	baseQuery := exerciseSubmissionRepo.DB.Model(models.ExerciseSubmission{})
	if groupId != nil {
		baseQuery.Joins("JOIN students ON students.stu_id = exercise_submissions.stu_id").
			Where("students.group_id = ?", groupId)
	}
	baseQuery.Count(&total)
	return total
}
