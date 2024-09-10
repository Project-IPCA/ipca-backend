package requests

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type AddStudentsTextRequest struct {
	StudentsData string    `json:"students_data" validate:"required"       example:"1 64010090 A B"`
	GroupID      uuid.UUID `json:"group_id"      validate:"required,uuid4" example:"00000000-0000-0000-0000-000000000000"`
}

func (asr AddStudentsTextRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&asr)
}
