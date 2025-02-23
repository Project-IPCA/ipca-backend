package models

import (
	"github.com/google/uuid"
)

type LabClassInfo struct {
	ChapterID               uuid.UUID                `gorm:"column:chapter_id;type:varchar(36);primaryKey"`
	ChapterIndex            int                      `gorm:"column:chapter_index;not null"`
	Name                    string                   `gorm:"column:name;type:varchar(256);not null"`
	FullMark                int                      `gorm:"column:fullmark;not null"`
	NoItems                 int                      `gorm:"column:no_items;not null;default:5"`
	Language                *string                  `gorm:"type:enum('PYTHON','C');default:'PYTHON'"`
	LabExercises            []LabExercise            `gorm:"foreignKey:ChapterID"`
	GroupChapterPermissions []GroupChapterPermission `gorm:"foreignKey:ChapterID"`
}

func (LabClassInfo) TableName() string {
	return "lab_class_infos"
}
