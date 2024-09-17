package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type StudentRepositoryQ interface {
	GetUserByKmitlID(
		student *models.Student,
		kmitlId string,
	)
	GetStudentByStuID(
		student *models.Student,
		stuId uuid.UUID,
	)
}

type StudentRepository struct {
	DB *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{DB: db}
}

func (studentRepository *StudentRepository) GetUserByKmitlID(
	student *models.Student,
	kmitlId string,
) {
	studentRepository.DB.Where("kmitl_id = ?", kmitlId).Find(student)
}

func (studentRepository *StudentRepository) GetStudentByStuID(
	student *models.Student,
	stuId uuid.UUID,
) {
	studentRepository.DB.Where("stu_id = ?", stuId).Find(student)
}

func (studentRepository *StudentRepository) GetStudentAssignChapter(studentAssignChapterItems *[]models.StudentAssignmentChapterItem,stuId uuid.UUID,chapterId uuid.UUID){
	studentRepository.DB.Where("stu_id = ? AND chapter_id = ?",stuId,chapterId).Order("item_id ASC").Find(studentAssignChapterItems)
}

func (studentRepository *StudentRepository) GetAllStudentAssignChapter(studentAssignChapterItems *[]models.StudentAssignmentChapterItem,stuId uuid.UUID){
	studentRepository.DB.Where("stu_id = ?",stuId).Find(studentAssignChapterItems)
}