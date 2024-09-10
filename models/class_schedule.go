package models

import (
	"github.com/google/uuid"
)

type ClassSchedule struct {
	GroupID                 uuid.UUID                `gorm:"type:varchar(36);primaryKey"`
	Number                  *int                     `gorm:"default:null"`
	Name                    string                   `gorm:"type:varchar(80)"`
	DeptID                  uuid.UUID                `gorm:"type:varchar(36)"`
	Department              Department               `gorm:"foreignKey:DeptID;references:DeptID"`
	SupervisorID            *uuid.UUID               `gorm:"type:varchar(36);default:null"`
	Supervisor              Supervisor               `gorm:"foreignKey:SupervisorID;references:SupervisorID"`
	Day                     *string                  `gorm:"type:enum('SUNDAY','MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY','SATURDAY');default:null"`
	TimeStart               *string                  `gorm:"type:time;default:null"`
	TimeEnd                 *string                  `gorm:"type:time;default:null"`
	Year                    *int                     `gorm:"default:null"`
	Semester                *int                     `gorm:"default:null"`
	AllowUploadPic          bool                     `gorm:"default:true"`
	AllowLogin              bool                     `gorm:"default:true"`
	AllowSubmit             bool                     `gorm:"default:true"`
	AllowExercise           bool                     `gorm:"default:true"`
	Students                []Student                `gorm:"foreignKey:GroupID"`
	ClassLabStaffs          []ClassLabStaff          `gorm:"foreignKey:ClassID"`
	GroupChapterPermissions []GroupChapterPermission `gorm:"foreignKey:ClassID"`
}

func (ClassSchedule) TableName() string {
	return "class_schedules"
}
