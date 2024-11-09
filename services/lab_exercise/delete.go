package labexercise

import (
	"fmt"
	"github.com/Project-IPCA/ipca-backend/models"
)

func (service *Service) Delete(
	labExercise *models.LabExercise,
) error {
	err := service.DB.Delete(*labExercise)
	if err.Error != nil {
		return fmt.Errorf("error while delete item : %v", err.Error)
	}
	return nil
}
