package groupassignmentchapteritem

import (
	"gorm.io/gorm"
)

type ServiceWrapper interface{}

type Service struct {
	DB *gorm.DB
}

func NewGroupAssignmentChapterItemService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
