package activitylog

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceWrapper interface {
	Create(
		groupId *uuid.UUID,
		username string,
		remoteIp string,
		remotePort *int,
		agent *string,
		pageName string,
		action string,
		ci *uint,
	) error
}

type Service struct {
	DB *gorm.DB
}

func NewActivityLogService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
