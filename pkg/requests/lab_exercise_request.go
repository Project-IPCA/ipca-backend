package requests

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type BasicKeywordConstrains struct {
	UserDefinedConstraints string `json:"user_defined_constraints"`
	SuggestedConstraints   string `json:"suggested_constraints"`
}

type CreateLabExerciseRequest struct {
	ChapterID          *uuid.UUID             `json:"chapter_id"          validate:"required" example:"00000000-0000-0000-0000-000000000000"`
	Level              string                 `json:"level"               validate:"required" example:"1"`
	Name               string                 `json:"name"                validate:"required" example:"Pyramid"`
	Content            string                 `json:"content"             validate:"required" example:"lorem"`
	Sourcecode         string                 `json:"sourcecode"          validate:"required" example:"lorem"`
	KeywordConstraints BasicKeywordConstrains `json:"keyword_constraints"`
}

func (cle CreateLabExerciseRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}
