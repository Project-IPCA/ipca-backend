package requests

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SuggestedConstraintsData struct {
	Limit   int    `json:"limit"`
	Keyword string `json:"keyword"`
}

type UserDefinedConstraintsData struct {
	Active  bool   `json:"active"`
	Keyword string `json:"keyword"`
	Type    string `json:"type"`
	Limit   int    `json:"limit"`
}

type UserDefinedConstraintsList struct {
	Classes      []UserDefinedConstraintsData `json:"classes"`
	Functions    []UserDefinedConstraintsData `json:"functions"`
	Imports      []UserDefinedConstraintsData `json:"imports"`
	Methods      []UserDefinedConstraintsData `json:"methods"`
	ReverseWords []UserDefinedConstraintsData `json:"reserved_words"`
	Variables    []UserDefinedConstraintsData `json:"variables"`
}

type SuggestedConstraintsList struct {
	Classes      []SuggestedConstraintsData `json:"classes"`
	Functions    []SuggestedConstraintsData `json:"functions"`
	Imports      []SuggestedConstraintsData `json:"imports"`
	Methods      []SuggestedConstraintsData `json:"methods"`
	ReverseWords []SuggestedConstraintsData `json:"reserved_words"`
	Variables    []SuggestedConstraintsData `json:"variables"`
}

type KeywordConstrains struct {
	UserDefinedConstraints UserDefinedConstraintsList `json:"user_defined_constraints"`
	SuggestedConstraints   SuggestedConstraintsList   `json:"suggested_constraints"`
}

type CreateLabExerciseRequest struct {
	ChapterID          *uuid.UUID        `json:"chapter_id"          validate:"required" example:"00000000-0000-0000-0000-000000000000"`
	Level              string            `json:"level"               validate:"required" example:"1"`
	Name               string            `json:"name"                validate:"required" example:"Pyramid"`
	Content            string            `json:"content"             validate:"required" example:"lorem"`
	Sourcecode         string            `json:"sourcecode"          validate:"required" example:"lorem"`
	KeywordConstraints KeywordConstrains `json:"keyword_constraints"`
}

func (cle CreateLabExerciseRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}

type UpdateLabExerciseRequest struct {
	JobID              *uuid.UUID        `json:"job_id"          validate:"required" example:"00000000-0000-0000-0000-000000000000"`
	ExerciseID         *uuid.UUID        `json:"exercise_id"          validate:"required" example:"00000000-0000-0000-0000-000000000000"`
	Name               string            `json:"name"                validate:"required" example:"Pyramid"`
	Content            string            `json:"content"             validate:"required" example:"lorem"`
	Sourcecode         string            `json:"sourcecode"          validate:"required" example:"lorem"`
	KeywordConstraints KeywordConstrains `json:"keyword_constraints"`
}

func (cle UpdateLabExerciseRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}
