package models

import (
	"github.com/google/uuid"
)

type Student struct {
	StuID       uuid.UUID                      `gorm:"primaryKey;type:varchar(36)"`
	KmitlID     string                         `gorm:"not null;type:varchar(10)"`
	GroupID     *uuid.UUID                     `gorm:"type:varchar(36);default:null"`
	Note        *string                        `gorm:"type:varchar(64);default:null"`
	MidCore     float64                        `gorm:"not null;default:0"`
	CanSubmit   bool                           `gorm:"type:boolean;not null;default:true"`
	User        *User                          `gorm:"foreignKey:StuID;references:UserID"`
	Group       *ClassSchedule                 `gorm:"foreignKey:GroupID;references:GroupID"`
	Assignments []StudentAssignmentChapterItem `gorm:"foreignKey:StuID;"`
}

func (Student) TableName() string {
	return "students"
}

type StudentWithAggregate struct {
	Student
	TotalMarks int `gorm:"column:total_marks"`
}
