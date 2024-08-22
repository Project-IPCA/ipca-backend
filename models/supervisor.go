package models

import (
	"github.com/google/uuid"
)

type Supervisor struct {
	SupervisorID   uuid.UUID       `gorm:"primaryKey;type:varchar(36)"`
	Dept           *string         `gorm:"type:varchar(40);default:null"`
	ClassLabStaffs []ClassLabStaff `gorm:"foreignKey:StaffID"`
}

func (Supervisor) TableName() string {
	return "supervisors"
}
