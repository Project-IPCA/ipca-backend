package models

import (
	"github.com/google/uuid"
)

type Supervisor struct {
	SupervisorID   uuid.UUID       `gorm:"primaryKey;type:varchar(36)"`
	User           User            `gorm:"foreignKey:SupervisorID;references:UserID"`
	ClassSchedules []ClassSchedule `gorm:"foreignKey:SupervisorID"`
	ClassLabStaffs []ClassLabStaff `gorm:"foreignKey:StaffID"`
}

func (Supervisor) TableName() string {
	return "supervisors"
}
