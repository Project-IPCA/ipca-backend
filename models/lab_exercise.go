package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type LabExercise struct {
	ExerciseID             uuid.UUID          `gorm:"type:varchar(36);primaryKey;column:exercise_id"`
	ChapterID              *uuid.UUID         `gorm:"type:varchar(36);column:chapter_id"`
	Level                  *string            `gorm:"type:enum('0','1','2','3','4','5','6')"`
	Name                   *string            `gorm:"type:varchar(1024)"`
	Content                *string            `gorm:"type:mediumtext"`
	Testcase               string             `gorm:"type:enum('NO_INPUT','YES','UNDEFINED');not null;default:'NO_INPUT'"`
	Sourcecode             *string            `gorm:"type:varchar(50)"`
	FullMark               int                `gorm:"type:int;not null;default:10"`
	AddedDate              time.Time          `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	LastUpdate             *time.Time         `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UserDefinedConstraints *json.RawMessage   `gorm:"type:json;column:user_defined_constraints"`
	SuggestedConstraints   *json.RawMessage   `gorm:"type:json;column:suggested_constraints"`
	AddedBy                *string            `gorm:"type:varchar(40)"`
	CreatedBy              *uuid.UUID         `gorm:"type:varchar(36)"`
	Chapter                LabClassInfo       `gorm:"foreignKey:ChapterID;references:ChapterID"`
	TestcaseList           []ExerciseTestcase `gorm:"foreignKey:ExerciseID;references:ExerciseID"`
}

func (LabExercise) TableName() string {
	return "lab_exercises"
}
