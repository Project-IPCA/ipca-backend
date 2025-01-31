package models

import (
	"github.com/google/uuid"
)

type ClassLabStaff struct {
	ClassID       uuid.UUID     `gorm:"column:class_id;type:varchar(36);primaryKey;not null"`
	StaffID       uuid.UUID     `gorm:"column:staff_id;type:varchar(36);primaryKey;not null"`
	Role          string        `gorm:"column:role;type:enum('EXECUTIVE','SUPERVISOR','TA','BEYONDER')"`
	User          User          `gorm:"foreignKey:StaffID;references:UserID"`
	ClassSchedule ClassSchedule `gorm:"foreignKey:ClassID;references:GroupID"`
}

func (ClassLabStaff) TableName() string {
	return "class_lab_staffs"
}
