package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID     uuid.UUID   `gorm:"type:varchar(36);primary_key"`
	Username   string      `gorm:"type:varchar(30);unique;not null"`
	Password   string      `gorm:"type:varchar(60)"`
	FirstName  *string     `gorm:"column:f_name;type:varchar(10)"`
	LastName   *string     `gorm:"column:l_name;type:varchar(32)"`
	Nickname   *string     `gorm:"type:varchar(50)"`
	Gender     *string     `gorm:"type:enum('MALE','FEMALE','OTHER')"`
	DOB        *time.Time  `gorm:"column:dob"`
	Avatar     *string     `gorm:"type:varchar(128)"`
	Role       *string     `gorm:"type:enum('ADMIN','EDITOR','AUTHOR','STUDENT','SUPERVISOR','STAFF','TA')"`
	Email      *string     `gorm:"type:varchar(64)"`
	Tel        *string     `gorm:"type:varchar(10)"`
	Added      time.Time   `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	LastLogin  *time.Time  `gorm:"column:last_login"`
	LastSeen   time.Time   `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	IsOnline   bool        `gorm:"type:boolean;not null;default:false"`
	IsActive   bool        `gorm:"type:boolean;not null;default:true"`
	AddedBy    *string     `gorm:"type:varchar(40)"`
	CISession  *int        `gorm:"column:ci_session"`
	SessionID  *string     `gorm:"column:session_id;type:varchar(50)"`
	Supervisor *Supervisor `gorm:"foreignKey:SupervisorID"`
	Student    *Student    `gorm:"foreignKey:StuID"`
	TA         *TA         `gorm:"foreignKey:TaID"`
}

func (User) TableName() string {
	return "users"
}
