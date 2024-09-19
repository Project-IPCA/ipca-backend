package repositories

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupChapterSelectedItemRepository struct {
	DB *gorm.DB
}

func NewGroupChapterSelectedItemRepository(db *gorm.DB) *GroupChapterSelectedItemRepository {
	return &GroupChapterSelectedItemRepository{DB: db}
}

func (groupChapterSelectedItemRepo *GroupChapterSelectedItemRepository) GetSelectedItemByGroupChapterItemId(
	selectedItems *[]models.GroupChapterSelectedItem,
	groupId uuid.UUID,
	chapterId uuid.UUID,
	itemId	int,
){
	groupChapterSelectedItemRepo.DB.
	Preload("ClassSchedule").
    Preload("LabClassInfo").
    Preload("LabExercise").
	Where("group_id = ? AND chapter_id = ? AND item_id = ?",groupId,chapterId,itemId).Find(selectedItems)
}