package models

import (
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	LogID      string     `gorm:"type:varchar(26);primary_key;column:log_id"`
	Timestamp  time.Time  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;column:timestamp"`
	GroupID    *uuid.UUID `gorm:"type:varchar(36);column:group_id"`
	Username   string     `gorm:"type:varchar(30);not null;column:username"`
	RemoteIP   string     `gorm:"type:varchar(15);not null;column:remote_ip"`
	RemotePort *int       `gorm:"column:remote_port"`
	Agent      *string    `gorm:"type:varchar(255);column:agent"`
	PageName   string     `gorm:"type:varchar(25);not null;column:page_name"`
	Action     string     `gorm:"type:text;not null;column:action"`
	CI         *uint      `gorm:"column:ci"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}