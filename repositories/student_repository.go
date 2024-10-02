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

type StudentWithAssignments struct {
	Username  string
	FirstName *string
	LastName  *string
	Assignments []AssignmentScore
}

type AssignmentScore struct {
	ChapterID uuid.UUID
	ItemID    int
	ExerciseID *uuid.UUID
	Marking    int
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

func (studentRepository *StudentRepository) GetStudentsAndAssignmentScoreByGroupID(
	student *[]models.Student,
	groupId uuid.UUID,
) {
	studentRepository.DB.Preload("User").Preload("Assignments").Where("group_id = ?", groupId).Find(student)
}