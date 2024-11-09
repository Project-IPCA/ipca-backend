package repositories

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentAssignChapterItemRepository struct {
	DB *gorm.DB
}

func NewStudentAssignChapterItemRepository(db *gorm.DB) *StudentAssignChapterItemRepository {
	return &StudentAssignChapterItemRepository{DB: db}
}

func (studentAssignChapterItemRepo *StudentAssignChapterItemRepository) GetStudentAssignChapter(studentAssignChapterItems *[]models.StudentAssignmentChapterItem, stuId uuid.UUID, chapterId uuid.UUID) {
	studentAssignChapterItemRepo.DB.Where("stu_id = ? AND chapter_id = ?", stuId, chapterId).Order("item_id ASC").Find(studentAssignChapterItems)
}

func (studentAssignChapterItemRepo *StudentAssignChapterItemRepository) GetAllStudentAssignChapter(studentAssignChapterItems *[]models.StudentAssignmentChapterItem, stuId uuid.UUID) {
	studentAssignChapterItemRepo.DB.Where("stu_id = ?", stuId).Find(studentAssignChapterItems)
}

func (studentAssignChapterItemRepo *StudentAssignChapterItemRepository) GetStudentAssignChapterItem(studentAssignChapterItems *models.StudentAssignmentChapterItem, stuId uuid.UUID, chapterId uuid.UUID, itemId int) error {
	err := studentAssignChapterItemRepo.DB.Where("stu_id = ? AND chapter_id = ? AND item_id = ?", stuId, chapterId, itemId).First(studentAssignChapterItems)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
