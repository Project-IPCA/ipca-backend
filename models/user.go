package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	StudentID string    `gorm:"type:char(10)"            json:"student_id"`
	FirstName string    `                                json:"f_name"`
	LastName  string    `                                json:"l_name"`
}
