package requests

import (
	"github.com/go-playground/validator/v10"
)

type SaveExerciseTestcaseRequest struct {
	ExerciseID string `json:"exercise_id" validate:"required"`
	JobID		string `json:"job_id" validate:"required"`
	RemoveList 	[]string `json:"removed_list" validate:"required"`
	TestCaseList []ExerciseTestcaseReq `json:"testcase_list" validate:"required"`
}

func (cle SaveExerciseTestcaseRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}

type AddTestcaseRabitMessage struct{
	JobId string `json:"job_id"`
	JobType string `json:"job_type"`
	ExerciseId string `json:"exercise_id"`
	TestcaseList []ExerciseTestcaseReq `json:"testcase_list"`
	Sourcecode string `json:"sourcecode"`
}

type ExerciseTestcaseReq struct {
	TestcaseID       *string `json:"testcase_id"`
	ExerciseID       string `json:"exercise_id"`
	IsReady          string `json:"is_ready"`
	TestcaseContent  string `json:"testcase_content"`
	IsActive         bool   `json:"is_active"`
	IsShowStudent    bool   `json:"show_to_student"`
	TestcaseNote     string `json:"testcase_note"`
	TestcaseOutput   string `json:"testcase_output"`
	TestcaseError    string `json:"testcase_error"`
}