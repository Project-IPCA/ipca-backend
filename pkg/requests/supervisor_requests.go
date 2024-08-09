package requests

import (
	"github.com/go-playground/validator/v10"
)

type AddStudentsTextRequest struct {
	StudentsData string `json:"students_data" validate:"required" example:"1 64010090 A B"`
}

func (asr AddStudentsTextRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&asr)
}
