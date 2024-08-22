package models

import (
	"github.com/google/uuid"
)

type GroupAssignmentExercise struct {
	GroupID    uuid.UUID `gorm:"type:varchar(36);primaryKey;not null"`
	ExerciseID uuid.UUID `gorm:"type:varchar(36);primaryKey;not null"`
	Selected   bool      `gorm:"not null;default:true"`

	Group    ClassSchedule `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Exercise LabExercise   `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

func (GroupAssignmentExercise) TableName() string {
	return "group_assignment_exercises"
}
