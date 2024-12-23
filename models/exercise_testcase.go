package models

import (
	"github.com/google/uuid"
)

type ExerciseTestcase struct {
	TestcaseID       *uuid.UUID `gorm:"column:testcase_id;type:varchar(36);primaryKey" json:"testcase_id"`
	ExerciseID       uuid.UUID `gorm:"column:exercise_id;type:varchar(36);not null" json:"excerise_id"`
	IsReady          string `gorm:"column:is_ready;type:varchar(3);not null;default:'yes'" json:"is_ready"`
	TestcaseContent  string `gorm:"column:testcase_content;type:varchar(1024);not null" json:"testcase_content"`
	IsActive         *bool   `gorm:"column:is_active;type:tinyint(1);default:1" json:"is_active"`
	IsShowStudent    *bool   `gorm:"column:is_show_student;type:tinyint(1);default:1" json:"show_to_student"`
	TestcaseNote     *string `gorm:"column:testcase_note;type:varchar(1024)" json:"testcase_note"`
	TestcaseOutput   *string `gorm:"column:testcase_output;type:mediumtext" json:"testcase_output"`
	TestcaseError    *string `gorm:"column:testcase_error;type:varchar(4096)" json:"testcase_error"`
}

func (ExerciseTestcase) TableName() string {
	return "exercise_testcases"
}