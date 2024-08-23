package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type ClassScheduleRepositoryQ interface {
	GetClassScheduleByGroupID(
		classSchedule *models.ClassSchedule,
		groupId uuid.UUID,
	)
}

type ClassScheduleRepository struct {
	DB *gorm.DB
}

func NewClassScheduleRepository(db *gorm.DB) *ClassScheduleRepository {
	return &ClassScheduleRepository{DB: db}
}

func (classScheduleRepository *ClassScheduleRepository) GetClassScheduleByGroupID(
	classSchedule *models.ClassSchedule,
	groupId uuid.UUID,
) {
	classScheduleRepository.DB.Where("group_id = ?", groupId).Find(classSchedule)
}

func (classScheduleRepository *ClassScheduleRepository) GetClassScheduleByNumber(
	classSchedule *models.ClassSchedule,
	number int,
) {
	classScheduleRepository.DB.Where("number = ?", number).Find(classSchedule)
}
