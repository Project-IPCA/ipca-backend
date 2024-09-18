package groupchapterselecteditem

import "gorm.io/gorm"

type Service struct {
	DB *gorm.DB
}

func NewGroupChapterSelectedItemService(db *gorm.DB) *Service {
	return &Service{DB: db}
}