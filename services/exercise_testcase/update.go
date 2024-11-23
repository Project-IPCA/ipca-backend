package exercisetestcase

import (
	"fmt"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
)

func (exerciseTestcaseService *Service) UpdateTestcaseIsReadyByExerciseID(
	labExerciseId uuid.UUID,
	isReady string,
) error {
	err := exerciseTestcaseService.DB.Model(models.ExerciseTestcase{}).Where("exercise_id = ?", labExerciseId).Update("is_ready", isReady)
	if err.Error != nil {
		return fmt.Errorf("failed to update exercise testcases: %w", err.Error)
	}
	return nil
}
