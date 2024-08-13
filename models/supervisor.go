package models

import (
	"github.com/google/uuid"
)

type Supervisor struct {
	SupervisorID uuid.UUID `gorm:"primaryKey;type:varchar(36)"`
	Deptment     *string   `gorm:"type:varchar(40);default:null"`
}

func (Supervisor) TableName() string {
	return "supervisors"
}
