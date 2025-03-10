package models

import (
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	LogID      string     `json:"log_id" gorm:"type:varchar(26);primary_key;column:log_id"`
	Timestamp  time.Time  `json:"timestamp" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;column:timestamp"`
	GroupID    *uuid.UUID `json:"group_id" gorm:"type:varchar(36);column:group_id"`
	Username   string     `json:"username" gorm:"type:varchar(30);not null;column:username"`
	RemoteIP   string     `json:"remote_ip" gorm:"type:varchar(15);not null;column:remote_ip"`
	RemotePort *int       `json:"remote_port" gorm:"column:remote_port"`
	Agent      *string    `json:"agent" gorm:"type:varchar(255);column:agent"`
	PageName   string     `json:"page_name" gorm:"type:varchar(25);not null;column:page_name"`
	Action     string     `json:"action" gorm:"type:text;not null;column:action"`
	CI         *uint      `json:"ci" gorm:"column:ci"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}

type LogExerciseSubmissionAction struct {
	StuId              uuid.UUID `json:"stu_id"`
	JobId              uuid.UUID `json:"job_id"`
	Status             string    `json:"status"`
	SubmissionId       uuid.UUID `json:"submission_id"`
	Attempt            string    `json:"attempt"`
	SourcecodeFilename string    `json:"sourcecode_filename"`
}
