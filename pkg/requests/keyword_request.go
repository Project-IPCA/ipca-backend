package requests

import (
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/go-playground/validator/v10"
)

type GetKeywordListRequest struct {
	Sourcecode string `json:"sourcecode"          validate:"required" example:"lorem"`
}

func (cle GetKeywordListRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}

type PythonCheckKeywordRequest struct {
	Sourcecode          string                               `json:"sourcecode"          validate:"required" example:"lorem"`
	ExerciseKeywordList constants.PythonCheckKeywordCategory `json:"exercise_kw_list"          validate:"required"`
}

func (cle PythonCheckKeywordRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}

type CCheckKeywordRequest struct {
	Sourcecode          string                          `json:"sourcecode"          validate:"required" example:"lorem"`
	ExerciseKeywordList constants.CCheckKeywordCategory `json:"exercise_kw_list"          validate:"required"`
}

func (cle CCheckKeywordRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}
