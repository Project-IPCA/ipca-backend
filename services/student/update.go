package student

import (
	"github.com/Project-IPCA/ipca-backend/models"
)

func (userStudentService *Service) UpdateCanSubmit(
	student *models.Student,
	canSubmit bool,
) error {
	student.CanSubmit = canSubmit
	userStudentService.DB.Save(student)
	return nil
}
