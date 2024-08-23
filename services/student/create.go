package student

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (userStudentService *Service) Create(
	userId uuid.UUID,
	KmitlId string,
	groupId *uuid.UUID,
) error {
	userStudent := builders.NewUserStudentBuilder().
		SetStuID(userId).
		SetKmitlID(KmitlId).
		SetGroupID(*groupId).
		Build()
	userStudentService.DB.Create(&userStudent)
	return nil
}
