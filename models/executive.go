package models

import (
	"github.com/google/uuid"
)

type Executive struct {
	ExecutiveID uuid.UUID `gorm:"primaryKey;type:varchar(36)"`
	User        User      `gorm:"foreignKey:SupervisorID;references:UserID"`
}

func (Executive) TableName() string {
	return "executives"
}
