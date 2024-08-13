package models

import (
	"time"

	"github.com/google/uuid"
)

type ClassSchedule struct {
	GroupID        uuid.UUID  `gorm:"type:varchar(36);primaryKey"`
	Number         *int       `gorm:"default:null"`
	Name           string     `gorm:"type:varchar(80)"`
	DeptID         uuid.UUID  `gorm:"type:varchar(36)"`
	SupervisorID   *uuid.UUID `gorm:"type:varchar(36);default:null"`
	Day            *string    `gorm:"type:enum('SUNDAY','MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY','SATURDAY');default:null"`
	TimeStart      *time.Time `gorm:"type:time;default:null"`
	TimeEnd        *time.Time `gorm:"type:time;default:null"`
	Year           *int       `gorm:"default:null"`
	Semester       *int       `gorm:"default:null"`
	AllowUploadPic bool       `gorm:"default:true"`
	AllowLogin     bool       `gorm:"default:true"`
	AllowSubmit    bool       `gorm:"default:true"`
	AllowExercise  bool       `gorm:"default:true"`
	Supervisor     Supervisor `gorm:"foreignKey:SupervisorID"`
	Department     Department `gorm:"foreignKey:DeptID"`
}

func (ClassSchedule) TableName() string {
	return "class_schedules"
}
