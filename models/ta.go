package models

import (
	"github.com/google/uuid"
)

type TA struct {
	TaID    string     `gorm:"primaryKey;type:varchar(36)"`
	GroupID *uuid.UUID `gorm:"type:varchar(36);default:null"`
	Note    *string    `gorm:"type:varchar(64);default:null"`
}

func (TA) TableName() string {
	return "tas"
}
