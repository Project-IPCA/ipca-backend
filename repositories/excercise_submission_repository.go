package repositories

import (
	"strings"
	"time"

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
	groupId string, year string,language string,
) int64 {
	var total int64
	baseQuery := exerciseSubmissionRepo.DB.Model(models.ExerciseSubmission{})

	if groupId != "" || year != "" {
		baseQuery = baseQuery.Joins(
			"JOIN students ON students.stu_id = exercise_submissions.stu_id",
		)
	}

	if groupId != "" {
		baseQuery = baseQuery.Where("students.group_id = ?", groupId)
	}

	if year != "" {
		baseQuery = baseQuery.Joins("JOIN class_schedules ON class_schedules.group_id = students.group_id").
			Where("class_schedules.year = ?", year)
	}

	if language != ""{
		baseQuery = baseQuery.Where("exercise_submissions.language = ?",strings.ToUpper(language))
	}

	baseQuery.Count(&total)

	return total
}

func (exerciseSubmissionRepo *ExerciseSubmissionRepository) GetSubmissionsByDate(
	submissions *int64,
	groupId string, year string,
	date time.Time,
	language string,
) {
	baseQuery := exerciseSubmissionRepo.DB.Model(models.ExerciseSubmission{})

	if groupId != "" || year != "" {
		baseQuery = baseQuery.Joins(
			"JOIN students ON students.stu_id = exercise_submissions.stu_id",
		)
	}

	if groupId != "" {
		baseQuery = baseQuery.Where("students.group_id = ?", groupId)
	}

	if year != "" {
		baseQuery = baseQuery.Joins("JOIN class_schedules ON class_schedules.group_id = students.group_id").
			Where("class_schedules.year = ?", year)
	}

	if language != ""{
		baseQuery = baseQuery.Where("exercise_submissions.language = ?",strings.ToUpper(language))
	}

	baseQuery = baseQuery.Where(
		"DATE(exercise_submissions.time_submit) = ?",
		date.Format("2006-01-02"),
	)

	baseQuery.Count(submissions)
}
