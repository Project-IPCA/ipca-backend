package labexercise

import (
	"encoding/json"
	"fmt"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/requests"
)

func (labExerciseService *Service) UpdateLabExercise(
	labExercise *models.LabExercise,
	updateData requests.UpdateLabExerciseRequest,
) error {
	suggestedJson, err := json.Marshal(updateData.KeywordConstraints.SuggestedConstraints)
	if err != nil {
		return fmt.Errorf("error while marshal suggestedConstraint : %v", err.Error())
	}
	userConstrainJson, err := json.Marshal(updateData.KeywordConstraints.UserDefinedConstraints)
	if err != nil {
		return fmt.Errorf("error while marshal userConstraint : %v", err.Error())
	}
	rawSuggested := json.RawMessage(suggestedJson)
	rawUserConstrain := json.RawMessage(userConstrainJson)
	labExercise.Name = &updateData.Name
	labExercise.Content = &updateData.Content
	labExercise.SuggestedConstraints = &rawSuggested
	labExercise.UserDefinedConstraints = &rawUserConstrain

	labExerciseService.DB.Save(labExercise)
	return nil
}

func (labExerciseService *Service) UpdateLabExerciseSourcecode(
	labExercise *models.LabExercise,
	filename string,
) {
	labExercise.Sourcecode = &filename
	labExerciseService.DB.Save(labExercise)
}
