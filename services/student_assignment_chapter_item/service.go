package studentassignmentchapteritem

import (
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewStudentAssignmentChapterItem(db *gorm.DB) *Service {
	return &Service{DB: db}
}