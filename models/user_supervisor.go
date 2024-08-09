package models

import (
	"time"

	"github.com/google/uuid"
)

type UserSupervisor struct {
	SupervisorID uuid.UUID  `gorm:"primaryKey;type:varchar(36)"                       json:"supervisor_id"`
	FirstName    *string    `gorm:"type:varchar(50);default:null"                     json:"supervisor_firstname"`
	LastName     *string    `gorm:"type:varchar(50);default:null"                     json:"supervisor_lastname"`
	Nickname     *string    `gorm:"type:varchar(50);default:null"                     json:"supervisor_nickname"`
	Gender       *string    `gorm:"type:enum('MALE', 'FEMALE', 'OTHER');default:null" json:"supervisor_gender"`
	DOB          *time.Time `gorm:"type:date;default:null"                            json:"supervisor_dob"`
	Avatar       *string    `gorm:"type:varchar(64);default:null"                     json:"supervisor_avatar"`
	Email        *string    `gorm:"type:varchar(64);default:null"                     json:"supervisor_email"`
	Tel          *string    `gorm:"type:varchar(10);default:null"                     json:"supervisor_tel"`
	Department   *string    `gorm:"type:varchar(40);default:null"                     json:"supervisor_department"`
}

func (UserSupervisor) TableName() string {
	return "user_supervisor"
}
