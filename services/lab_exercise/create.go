package labexercise

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/pkg/constants"
	"github.com/Project-IPCA/ipca-backend/pkg/requests"
	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (labExerciseService *Service) CreatePythonExercise(
	request *requests.CreatePythonExerciseRequest,
	supervisorId *uuid.UUID,
	supervisorName string,
) (uuid.UUID, error) {
	exerciseId := uuid.New()
	labExercise := builders.NewLabExerciseBuilder().
		SetExerciseID(exerciseId).
		SetChapterID(request.ChapterID).
		SetLevel(&request.Level).
		SetName(&request.Name).
		SetContent(&request.Content).
		SetSourcecode(&request.Sourcecode).
		SetAddedBy(&supervisorName).
		SetCreatedBy(supervisorId).
		SetLanguage(&constants.ExerciseLanguage.Python).
		Build()

	labExerciseService.DB.Create(&labExercise)
	return exerciseId, nil
}

func (labExerciseService *Service) CreateCExercise(
	request *requests.CreateCExerciseRequest,
	supervisorId *uuid.UUID,
	supervisorName string,
) (uuid.UUID, error) {
	exerciseId := uuid.New()
	labExercise := builders.NewLabExerciseBuilder().
		SetExerciseID(exerciseId).
		SetChapterID(request.ChapterID).
		SetLevel(&request.Level).
		SetName(&request.Name).
		SetContent(&request.Content).
		SetSourcecode(&request.Sourcecode).
		SetAddedBy(&supervisorName).
		SetCreatedBy(supervisorId).
		SetLanguage(&constants.ExerciseLanguage.C).
		Build()

	labExerciseService.DB.Create(&labExercise)
	return exerciseId, nil
}

func (labExerciseService *Service) CreatePythonWithoutSourceCode(
	request *requests.CreatePythonExerciseRequest,
	supervisorId *uuid.UUID,
	supervisorName string,
) (uuid.UUID, error) {
	exerciseId := uuid.New()
	suggestedJson, err := json.Marshal(request.KeywordConstraints.SuggestedConstraints)
	if err != nil {
		return exerciseId, fmt.Errorf("error while marshal suggestedConstraint : %v", err.Error())
	}
	userConstrainJson, err := json.Marshal(request.KeywordConstraints.UserDefinedConstraints)
	if err != nil {
		return exerciseId, fmt.Errorf("error while marshal userConstraint : %v", err.Error())
	}
	rawSuggested := json.RawMessage(suggestedJson)
	rawUserConstrain := json.RawMessage(userConstrainJson)
	labExercise := builders.NewLabExerciseBuilder().
		SetExerciseID(exerciseId).
		SetChapterID(request.ChapterID).
		SetLevel(&request.Level).
		SetName(&request.Name).
		SetContent(&request.Content).
		SetAddedBy(&supervisorName).
		SetCreatedBy(supervisorId).
		SetSuggestedConstraints(&rawSuggested).
		SetUserDefinedConstraints(&rawUserConstrain).
		SetLanguage(&constants.ExerciseLanguage.Python).
		Build()

	labExerciseService.DB.Create(&labExercise)
	return exerciseId, nil
}

func (labExerciseService *Service) CreateCWithoutSourceCode(
	request *requests.CreateCExerciseRequest,
	supervisorId *uuid.UUID,
	supervisorName string,
) (uuid.UUID, error) {
	exerciseId := uuid.New()
	suggestedJson, err := json.Marshal(request.KeywordConstraints.SuggestedConstraints)
	if err != nil {
		return exerciseId, fmt.Errorf("error while marshal suggestedConstraint : %v", err.Error())
	}
	userConstrainJson, err := json.Marshal(request.KeywordConstraints.UserDefinedConstraints)
	if err != nil {
		return exerciseId, fmt.Errorf("error while marshal userConstraint : %v", err.Error())
	}
	rawSuggested := json.RawMessage(suggestedJson)
	rawUserConstrain := json.RawMessage(userConstrainJson)
	labExercise := builders.NewLabExerciseBuilder().
		SetExerciseID(exerciseId).
		SetChapterID(request.ChapterID).
		SetLevel(&request.Level).
		SetName(&request.Name).
		SetContent(&request.Content).
		SetAddedBy(&supervisorName).
		SetCreatedBy(supervisorId).
		SetSuggestedConstraints(&rawSuggested).
		SetUserDefinedConstraints(&rawUserConstrain).
		SetLanguage(&constants.ExerciseLanguage.C).
		Build()

	labExerciseService.DB.Create(&labExercise)
	return exerciseId, nil
}