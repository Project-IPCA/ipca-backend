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

type UserDefinedPythonConstraintsList struct {
	Classes      []UserDefinedConstraintsData `json:"classes"`
	Functions    []UserDefinedConstraintsData `json:"functions"`
	Imports      []UserDefinedConstraintsData `json:"imports"`
	Methods      []UserDefinedConstraintsData `json:"methods"`
	ReverseWords []UserDefinedConstraintsData `json:"reserved_words"`
	Variables    []UserDefinedConstraintsData `json:"variables"`
}

type SuggestedPythonConstraintsList struct {
	Classes      []SuggestedConstraintsData `json:"classes"`
	Functions    []SuggestedConstraintsData `json:"functions"`
	Imports      []SuggestedConstraintsData `json:"imports"`
	Methods      []SuggestedConstraintsData `json:"methods"`
	ReverseWords []SuggestedConstraintsData `json:"reserved_words"`
	Variables    []SuggestedConstraintsData `json:"variables"`
}

type PythonKeywordConstrains struct {
	UserDefinedConstraints UserDefinedPythonConstraintsList `json:"user_defined_constraints"`
	SuggestedConstraints   SuggestedPythonConstraintsList   `json:"suggested_constraints"`
}

type CreatePythonExerciseRequest struct {
	CreateLabExerciseBasic
	KeywordConstraints PythonKeywordConstrains `json:"keyword_constraints"`
}

func (cle CreatePythonExerciseRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}

type UserDefinedCConstraintsList struct {
	Functions    []UserDefinedConstraintsData `json:"functions"`
	Includes     []UserDefinedConstraintsData `json:"includes"`
	ReverseWords []UserDefinedConstraintsData `json:"reserved_words"`
	Variables    []UserDefinedConstraintsData `json:"variables"`
}

type SuggestedCConstraintsList struct {
	Functions    []SuggestedConstraintsData `json:"functions"`
	Includes     []SuggestedConstraintsData `json:"includes"`
	ReverseWords []SuggestedConstraintsData `json:"reserved_words"`
	Variables    []SuggestedConstraintsData `json:"variables"`
}

type CKeywordConstrains struct {
	UserDefinedConstraints UserDefinedCConstraintsList `json:"user_defined_constraints"`
	SuggestedConstraints   SuggestedCConstraintsList   `json:"suggested_constraints"`
}

type CreateLabExerciseBasic struct {
	ChapterID  *uuid.UUID `json:"chapter_id"          validate:"required" example:"00000000-0000-0000-0000-000000000000"`
	Level      string     `json:"level"               validate:"required" example:"1"`
	Name       string     `json:"name"                validate:"required" example:"Pyramid"`
	Content    string     `json:"content"             validate:"required" example:"lorem"`
	Sourcecode string     `json:"sourcecode"          validate:"required" example:"lorem"`
}

type CreateCExerciseRequest struct {
	CreateLabExerciseBasic
	KeywordConstraints CKeywordConstrains `json:"keyword_constraints"`
}

func (cle CreateCExerciseRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}

type UpdateExerciseBasic struct {
	JobID      *uuid.UUID `json:"job_id"          validate:"required" example:"00000000-0000-0000-0000-000000000000"`
	ExerciseID *uuid.UUID `json:"exercise_id"          validate:"required" example:"00000000-0000-0000-0000-000000000000"`
	Name       string     `json:"name"                validate:"required" example:"Pyramid"`
	Content    string     `json:"content"             validate:"required" example:"lorem"`
	Sourcecode string     `json:"sourcecode"          validate:"required" example:"lorem"`
}

type UpdatePythonExerciseRequest struct {
	UpdateExerciseBasic
	KeywordConstraints PythonKeywordConstrains `json:"keyword_constraints"`
}

func (cle UpdatePythonExerciseRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}

type UpdateCExerciseRequest struct {
	UpdateExerciseBasic
	KeywordConstraints CKeywordConstrains `json:"keyword_constraints"`
}

func (cle UpdateCExerciseRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(&cle)
}
