package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type GroupChapterPermissionRepositoryQ interface{}

type GroupChapterPermissionRepository struct {
	DB *gorm.DB
}

func NewGroupChapterPermissionRepository(db *gorm.DB) *GroupChapterPermissionRepository {
	return &GroupChapterPermissionRepository{DB: db}
}

func (groupChapterPermissionRepository *GroupChapterPermissionRepository) GetGroupChapterPermissionByPK(
	groupChapterPermission *models.GroupChapterPermission,
	classId uuid.UUID,
	chapterId uuid.UUID,
) {
	groupChapterPermissionRepository.DB.Where("class_id = ?", classId).
		Where("chapter_id = ?", chapterId).Find(groupChapterPermission)
}

func (groupChapterPermissionRepository *GroupChapterPermissionRepository) GetGroupChapterPermissionByGroupID(
	groupChapterPermission *[]models.GroupChapterPermission,
	classId uuid.UUID,
) {
	groupChapterPermissionRepository.DB.Preload("LabClassInfo").Where("class_id = ?", classId).Find(groupChapterPermission)
}
