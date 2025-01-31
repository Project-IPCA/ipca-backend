package models

import (
	"github.com/google/uuid"
)

type Executive struct {
	ExecutiveID uuid.UUID `gorm:"primaryKey;type:varchar(36)"`
	User        User      `gorm:"foreignKey:ExecutiveID;references:UserID"`
}

func (Executive) TableName() string {
	return "executives"
}
