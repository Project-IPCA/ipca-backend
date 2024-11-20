package exercisesubmission

import (
	"fmt"

	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
)

func (exerciseSubmissionService *Service) CancleSubmission(
	submission *models.ExerciseSubmission,
) error {
	submission.Status = constants.ExerciseStatus.Rejected
	submission.Marking = 0
	
	err := exerciseSubmissionService.DB.Save(submission)
	if err.Error != nil{
		return fmt.Errorf("error while update item : %v", err.Error)
	}

	return nil
}
