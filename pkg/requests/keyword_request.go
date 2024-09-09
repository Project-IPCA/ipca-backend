package requests

import (
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/go-playground/validator/v10"
)

type GetKeywordListRequest struct {
	Sourcecode         string                 `json:"sourcecode"          validate:"required" example:"lorem"`
}

type CheckKeywordRequest struct{
	Sourcecode         string                 `json:"sourcecode"          validate:"required" example:"lorem"`
	ExerciseKeywordList  constants.CheckKeywordCategory  `json:"exercise_kw_list"          validate:"required"`
}

func (cle GetKeywordListRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}

func (cle CheckKeywordRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}