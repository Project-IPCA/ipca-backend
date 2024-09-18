package models

import "github.com/google/uuid"

type GroupChapterSelectedItem struct {
	GroupID    uuid.UUID `gorm:"type:varchar(36);column:group_id;primaryKey"`
	ChapterID  uuid.UUID  `gorm:"type:varchar(36);column:chapter_id;primaryKey"`
	ItemID     int64      `gorm:"type:int;column:item_id"`
	ExerciseID uuid.UUID  `gorm:"type:varchar(36);column:exercise_id;primaryKey"`
	ClassSchedule *ClassSchedule `gorm:"foreignKey:ClassID;references:GroupID"`
	LabClassInfo  *LabClassInfo  `gorm:"foreignKey:ChapterID;references:ChapterID"`
	LabExercise	*LabExercise `gorm:"foreignKey:ExerciseID;references:ExerciseID"`
}

func (GroupChapterSelectedItem) TableName() string {
	return "group_chapter_selected_items"
}
