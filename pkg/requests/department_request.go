package requests

import "github.com/go-playground/validator/v10"

type CreateDepartmentRequest struct {
	Name    string `json:"name_th" validate:"required"`
	Name_EN string `json:"name_en" validate:"required"`
}

func (cle CreateDepartmentRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}
