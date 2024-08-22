package models

import (
	"github.com/google/uuid"
)

type ClassLabStaff struct {
	ClassID       uuid.UUID     `gorm:"column:class_id;type:varchar(36);primaryKey;not null"`
	StaffID       uuid.UUID     `gorm:"column:staff_id;type:varchar(36);primaryKey;not null"`
	Supervisor    Supervisor    `gorm:"foreignKey:StaffID;references:SupervisorID"`
	ClassSchedule ClassSchedule `gorm:"foreignKey:ClassID;references:GroupID"`
}

func (ClassLabStaff) TableName() string {
	return "class_lab_staffs"
}
