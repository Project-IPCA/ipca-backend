package requests

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SaveExerciseTestcaseRequest struct {
	ExerciseID   uuid.UUID             `json:"exercise_id" validate:"required"`
	JobID        uuid.UUID             `json:"job_id" validate:"required"`
	RemoveList   []string              `json:"removed_list" validate:"required"`
	TestCaseList []ExerciseTestcaseReq `json:"testcase_list" validate:"required"`
}

func (cle SaveExerciseTestcaseRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}

type AddTestcaseRabitMessage struct {
	JobId        uuid.UUID             `json:"job_id"`
	JobType      string                `json:"job_type"`
	ExerciseId   uuid.UUID             `json:"exercise_id"`
	TestcaseList []ExerciseTestcaseReq `json:"testcase_list"`
	Sourcecode   string                `json:"sourcecode"`
}

type ExerciseTestcaseReq struct {
	TestcaseID      *uuid.UUID `json:"testcase_id"`
	ExerciseID      uuid.UUID  `json:"exercise_id"`
	IsReady         string     `json:"is_ready"`
	TestcaseContent string     `json:"testcase_content"`
	IsActive        bool       `json:"is_active"`
	IsShowStudent   bool       `json:"show_to_student"`
	TestcaseNote    string     `json:"testcase_note"`
	TestcaseOutput  string     `json:"testcase_output"`
	TestcaseError   string     `json:"testcase_error"`
}