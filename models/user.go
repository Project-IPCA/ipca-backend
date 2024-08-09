package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID   `gorm:"primaryKey;type:varchar(36)"                                                                 json:"id"`
	UserStudent UserStudent `gorm:"foreignKey:StuID"`
	Username    string      `gorm:"unique;type:varchar(30);not null"                                                            json:"username"`
	Password    string      `gorm:"type:varchar(60)"                                                                            json:"password"`
	Role        string      `gorm:"type:enum('ADMIN', 'EDITOR', 'AUTHOR', 'STUDENT', 'SUPERVISOR', 'STAFF', 'TA');default:null" json:"role"`
	Added       time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP"                                                          json:"added"`
	LastLogin   *time.Time  `gorm:"default:null"                                                                                json:"last_login"`
	LastSeen    time.Time   `gorm:"not null;default:CURRENT_TIMESTAMP"                                                          json:"last_seen"`
	IsOnline    bool        `gorm:"not null;default:false"                                                                      json:"is_online"`
	IsActive    bool        `gorm:"not null;default:true"                                                                       json:"is_active"`
	AddedBy     *string     `gorm:"type:varchar(40);default:null"                                                               json:"added_by"`
	CiSession   *int        `gorm:"default:null"                                                                                json:"ci_session"`
	SessionID   *string     `gorm:"type:varchar(50);default:null"                                                               json:"session_id"`
}

func (User) TableName() string {
	return "user"
}
