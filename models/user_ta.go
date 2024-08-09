package models

import (
	"time"
)

type UserTA struct {
	TaID      string     `gorm:"primaryKey;type:varchar(36)"                       json:"ta_id"`
	Gender    *string    `gorm:"type:enum('MALE', 'FEMALE', 'OTHER');default:null" json:"ta_gender"`
	FirstName *string    `gorm:"type:varchar(40);default:null"                     json:"ta_firstname"`
	LastName  *string    `gorm:"type:varchar(32);default:null"                     json:"ta_lastname"`
	Nickname  *string    `gorm:"type:varchar(20);default:null"                     json:"ta_nickname"`
	DOB       *time.Time `gorm:"type:date;default:null"                            json:"ta_dob"`
	Avatar    *string    `gorm:"type:varchar(128);default:null"                    json:"ta_avatar"`
	Email     *string    `gorm:"type:varchar(64);default:null"                     json:"ta_email"`
	Tel       *string    `gorm:"type:varchar(10);default:null"                     json:"ta_tel"`
	Group     *int       `gorm:"default:null"                                      json:"ta_group"`
	Note      *string    `gorm:"type:varchar(64);default:null"                     json:"note"`
	DeptID    *int       `gorm:"default:null"                                      json:"ta_dept_id"`
}

func (UserTA) TableName() string {
	return "user_ta"
}
