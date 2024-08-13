package student

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (userStudentService *Service) Create(
	userId uuid.UUID,
	KmitlId string,
) error {
	userStudent := builders.NewUserStudentBuilder().
		SetStuID(userId).
		SetKmitlID(KmitlId).
		Build()
	userStudentService.DB.Create(&userStudent)
	return nil
}
