package requests

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type AddStudentsTextRequest struct {
	StudentsData string    `json:"students_data" validate:"required"       example:"1 64010090 A B"`
	GroupID      uuid.UUID `json:"group_id"      validate:"required,uuid4" example:"00000000-0000-0000-0000-000000000000"`
}

func (asr AddStudentsTextRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&asr)
}

type ExcerciseSubmitRequest struct {
	ChapterID  string    `json:"chapter_id" validate:"required"`
	ItemId     int       `json:"item_id"    validate:"required"`
	Sourcecode string    `json:"sourcecode" validate:"required"`
	JobId      uuid.UUID `json:"job_id"     validate:"required"`
}

func (esr ExcerciseSubmitRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&esr)
}

type ExerciseSubmissionRabbitMessage struct {
	JobId        uuid.UUID                 `json:"job_id"`
	JobType      string                    `json:"job_type"`
	LogData      LogDataInfo               `json:"log_data"`
	SubmissionId uuid.UUID                 `json:"submission_id"`
	SourceCode   string                    `json:"sourcecode"`
	TestCaseList []models.ExerciseTestcase `json:"testcase_list"`
	StudentID    uuid.UUID                 `json:"stu_id"`
	ChapterId    uuid.UUID                 `json:"chapter_id"`
	ItemId       int                       `json:"item_id"`
	Language     string                    `json:"language"`
}

type LogDataInfo struct {
	GroupID    uuid.UUID                          `json:"group_id"`
	Username   string                             `json:"username"`
	RemoteIP   string                             `json:"remote_ip"`
	RemotePort int                                `json:"remote_port"`
	Agent      string                             `json:"agent"`
	PageName   string                             `json:"page_name"`
	Action     models.LogExerciseSubmissionAction `json:"action"`
}

type CanSubmitRequest struct {
	CanSubmit bool `json:"can_submit" example:"true" validate:"omitempty"`
}

func (csr CanSubmitRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&csr)
}
