package models

import (
	"github.com/google/uuid"
)

type Department struct {
	DeptID         uuid.UUID       `gorm:"type:varchar(36);primaryKey"`
	Name           string          `gorm:"type:varchar(100);unique;not null;"`
	ClassSchedules []ClassSchedule `gorm:"foreignKey:DeptID"`
}

func (Department) TableName() string {
	return "departments"
}
