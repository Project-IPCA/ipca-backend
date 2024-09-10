package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type GroupAssignmentChapterItemRepositoryQ interface{}

type GroupAssignmentChapterItemRepository struct {
	DB *gorm.DB
}

func NewGroupAssignmentChapterItemRepository(db *gorm.DB) *GroupAssignmentChapterItemRepository {
	return &GroupAssignmentChapterItemRepository{DB: db}
}

func (groupAssignmentChapterItemRepository *GroupAssignmentChapterItemRepository) GetGroupAssignmentChapterItemsByPK(
	groupAssignmentChapterItem *models.GroupAssignmentChapterItem,
	groupId uuid.UUID,
	chapterId uuid.UUID,
	itemId int64,
) {
	groupAssignmentChapterItemRepository.DB.Where("group_id = ?", groupId).
		Where("chapter_id = ?", chapterId).
		Where("item_id = ?", itemId).
		Find(groupAssignmentChapterItem)
}
