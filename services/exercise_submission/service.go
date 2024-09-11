package exercisesubmission

import "gorm.io/gorm"

type Service struct {
	DB *gorm.DB
}

func NewExerciseSubmissionService(db *gorm.DB) *Service {
	return &Service{DB: db}
}