package classlabstaff

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (classLabStaffService *Service) Create(classId uuid.UUID, staffId uuid.UUID) error {
	classLabStaff := builders.NewClassLabStaffBuilder().
		SetClassID(classId).
		SetStaffID(staffId).
		Build()

	classLabStaffService.DB.Create(&classLabStaff)
	return nil
}
