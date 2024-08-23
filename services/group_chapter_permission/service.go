package groupchapterpermission

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceWrapper interface {
	Create(
		classId uuid.UUID,
		chapterId uuid.UUID,
		timeStart string,
		timeEnd string,
		allowSubmit bool,
		allowAccess bool,
		status string,
	) error
}

type Service struct {
	DB *gorm.DB
}

func NewGroupChapterPermissionService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
