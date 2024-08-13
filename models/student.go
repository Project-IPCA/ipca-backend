package models

import (
	"github.com/google/uuid"
)

type Student struct {
	StuID     uuid.UUID  `gorm:"primaryKey;type:varchar(36)"`
	KmitlID   string     `gorm:"not null;type:varchar(10)"`
	GroupID   *uuid.UUID `gorm:"type:varchar(36);default:null"`
	Note      *string    `gorm:"type:varchar(64);default:null"`
	DeptID    *uuid.UUID `gorm:"type:varchar(36);default:null"`
	MidCore   float64    `gorm:"not null;default:0"`
	CanSubmit bool       `gorm:"type:boolean;not null;default:true"`
}

func (Student) TableName() string {
	return "students"
}
