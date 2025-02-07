package repositories

import (
	// "time"

	// "github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type ActivityLogRepositoryQ interface{}

type ActivityLogRepository struct {
	DB *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) *ActivityLogRepository {
	return &ActivityLogRepository{DB: db}
}

func (activityLogRepo *ActivityLogRepository) GetActivityLogByGroupID(
	activityLog *[]models.ActivityLog,
	groupId string,
	lastTime string,
	limit int,
) (int64, error) {
	var total int64
	baseQuery := activityLogRepo.DB.Model(models.ActivityLog{}).Where("group_id = ?", groupId)

	baseQuery.Count(&total)

	if lastTime != "" {
		baseQuery = baseQuery.Where("timestamp < ?", lastTime)
	}

	err := baseQuery.Order("timestamp DESC").Limit(limit).Find(&activityLog)

	if err.Error != nil {
		return 0, err.Error
	}
	return total, nil
}
