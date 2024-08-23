package models

import (
	"time"

	"github.com/google/uuid"
)

type GroupChapterPermission struct {
	ClassID         uuid.UUID     `gorm:"column:class_id;type:varchar(36);primaryKey;not null"`
	ChapterID       uuid.UUID     `gorm:"column:chapter_id;type:varchar(36);primaryKey;not null"`
	AllowAccessType string        `gorm:"column:allow_access_type;type:enum('DENY','ALWAYS','TIMER','TIMER_PAUSED','DATETIME');not null;default:DENY"`
	AccessTimeStart *time.Time    `gorm:"column:access_time_start"`
	AccessTimeEnd   *time.Time    `gorm:"column:access_time_end"`
	AllowSubmitType string        `gorm:"column:allow_submit_type;type:enum('DENY','ALWAYS','TIMER','TIMER_PAUSED','DATETIME');not null;default:DENY"`
	SubmitTimeStart *time.Time    `gorm:"column:submit_time_start"`
	SubmitTimeEnd   *time.Time    `gorm:"column:submit_time_end"`
	AllowSubmit     bool          `gorm:"column:allow_submit;not null;default:true"`
	Status          string        `gorm:"column:status;type:enum('NA','READY','OPEN','CLOSE','STOP');not null;default:NA"`
	AllowAccess     bool          `gorm:"column:allow_access;not null;default:false"`
	TimeStart       *string       `gorm:"column:time_start;type:varchar(8)"`
	TimeEnd         *string       `gorm:"column:time_end;type:varchar(8)"`
	ClassSchedule   ClassSchedule `gorm:"foreignKey:ClassID;references:GroupID"`
}

func (GroupChapterPermission) TableName() string {
	return "group_chapter_permissions"
}
