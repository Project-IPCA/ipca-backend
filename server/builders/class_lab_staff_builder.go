package builders

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type ClassLabStaffBuilder struct {
	ClassID uuid.UUID
	StaffID uuid.UUID
	Role    string
}

func NewClassLabStaffBuilder() *ClassLabStaffBuilder {
	return &ClassLabStaffBuilder{}
}

func (classLabStaffBuilder *ClassLabStaffBuilder) SetClassID(
	classId uuid.UUID,
) (cls *ClassLabStaffBuilder) {
	classLabStaffBuilder.ClassID = classId
	return classLabStaffBuilder
}

func (classLabStaffBuilder *ClassLabStaffBuilder) SetStaffID(
	staffId uuid.UUID,
) (cls *ClassLabStaffBuilder) {
	classLabStaffBuilder.StaffID = staffId
	return classLabStaffBuilder
}

func (classLabStaffBuilder *ClassLabStaffBuilder) SetRole(
	role string,
) (cls *ClassLabStaffBuilder) {
	classLabStaffBuilder.Role = role
	return classLabStaffBuilder
}

func (classLabStaffBuilder *ClassLabStaffBuilder) Build() models.ClassLabStaff {
	classLabStaff := models.ClassLabStaff{
		ClassID: classLabStaffBuilder.ClassID,
		StaffID: classLabStaffBuilder.StaffID,
		Role:    classLabStaffBuilder.Role,
	}
	return classLabStaff
}
