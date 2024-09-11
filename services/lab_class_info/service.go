package labclassinfo
import (
	"gorm.io/gorm"
)

type ServiceWrapper interface {
	Create(
		chapterIndex int,
		name         string,
		fullMark     int,
		noItems      int,
	) error
}

type Service struct {
	DB *gorm.DB
}

func NewLabClassInfoService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
