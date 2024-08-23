package models

import (
	"github.com/google/uuid"
)

type LabClassInfo struct {
	ChapterID    uuid.UUID     `gorm:"type:varchar(36);primaryKey"`
	ChapterIndex int           `gorm:"not null"`
	Name         string        `gorm:"type:varchar(256);not null"`
	FullMark     int           `gorm:"not null"`
	NoItems      int           `gorm:"not null;default:5"`
	LabExercises []LabExercise `gorm:"foreignKey:ChapterID"`
}

func (LabClassInfo) TableName() string {
	return "lab_class_infos"
}
