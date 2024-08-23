package models

import (
	"github.com/google/uuid"
)

type GroupAssignmentChapterItem struct {
	GroupID        uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	ChapterID      uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	ItemID         int64     `gorm:"type:varchar(36);primaryKey"`
	ExerciseIDList *string   `gorm:"type:varchar(1024)"`
	FullMark       int       `gorm:"not null;default:2"`
	TimeStart      *string   `gorm:"default:null"`
	TimeEnd        *string   `gorm:"default:null"`
	Status         *string   `gorm:"type:enum('READY','CLOSED','STOP','OPEN');default:null"`
}

func (GroupAssignmentChapterItem) TableName() string {
	return "group_assignment_chapter_items"
}
