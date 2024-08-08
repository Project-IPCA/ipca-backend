package models

import (
	"time"

	"github.com/google/uuid"
)

type UserStudent struct {
	StuID        uuid.UUID  `gorm:"primaryKey;type:varchar(36)"                       json:"stu_id"`
	StuStuID     string     `gorm:"not null;type:varchar(10)"                         json:"stu_stu_id"`
	StuFirstname *string    `gorm:"type:varchar(40);default:null"                     json:"stu_firstname"`
	StuLastname  *string    `gorm:"type:varchar(32);default:null"                     json:"stu_lastname"`
	StuNickname  *string    `gorm:"type:varchar(20);default:null"                     json:"stu_nickname"`
	StuGender    *string    `gorm:"type:enum('MALE', 'FEMALE', 'OTHER');default:null" json:"stu_gender"`
	StuDOB       *time.Time `gorm:"type:date;default:null"                            json:"stu_dob"`
	StuAvatar    *string    `gorm:"type:varchar(128);default:null"                    json:"stu_avatar"`
	StuEmail     *string    `gorm:"type:varchar(64);default:null"                     json:"stu_email"`
	StuTel       *string    `gorm:"type:varchar(10);default:null"                     json:"stu_tel"`
	StuGroup     *int       `gorm:"default:null"                                      json:"stu_group"`
	Note         *string    `gorm:"type:varchar(64);default:null"                     json:"note"`
	StuDeptID    *int       `gorm:"default:null"                                      json:"stu_dept_id"`
	MidCore      float64    `gorm:"not null;default:0"                                json:"mid_core"`
	CanSubmit    string     `gorm:"type:varchar(3);not null;default:'YES'"            json:"can_submit"`
}

func (UserStudent) TableName() string {
	return "user_student"
}
