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

func (studentAssignChapterItemRepo *StudentAssignChapterItemRepository) GetAllStudentAssignChapterWithSubmission(studentAssignChapterItems *[]models.StudentAssignmentChapterItem, stuId uuid.UUID) {
	studentAssignChapterItemRepo.DB.Where("stu_id = ?", stuId).Preload("SubmissionList", "stu_id = ?", stuId).Find(studentAssignChapterItems)
}

func (studentAssignChapterItemRepo *StudentAssignChapterItemRepository) GetStudentAssignChapterItem(studentAssignChapterItems *models.StudentAssignmentChapterItem, stuId uuid.UUID, chapterId uuid.UUID, itemId int) error {
	err := studentAssignChapterItemRepo.DB.Where("stu_id = ? AND chapter_id = ? AND item_id = ?", stuId, chapterId, itemId).First(studentAssignChapterItems)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (studentAssignChapterItemRepo *StudentAssignChapterItemRepository) GetStudentChapterByChapterIDAndGroupOrYear(
	studentAssignChapterItems *[]models.StudentAssignmentChapterItem,
	groupId uuid.UUID,
	chapterId uuid.UUID,
	year string,
) error {
	baseQuery := studentAssignChapterItemRepo.DB.Model(models.StudentAssignmentChapterItem{}).
		Joins("JOIN lab_class_infos ON lab_class_infos.chapter_id = student_assignment_chapter_items.chapter_id").
		Joins("JOIN students ON students.stu_id = student_assignment_chapter_items.stu_id")

	if year != "" {
		baseQuery.Joins("JOIN class_schedules ON class_schedules.group_id = students.group_id").
			Where("class_schedules.year = ?", year)
	}

	if groupId != uuid.Nil {
		baseQuery.Where("students.group_id = ?", groupId)
	}

	err := baseQuery.Where("student_assignment_chapter_items.chapter_id = ?", chapterId).Order("lab_class_infos.chapter_index ASC").Find(studentAssignChapterItems)

	if err.Error != nil {
		return err.Error
	}
	return nil
}
