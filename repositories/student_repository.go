package repositories

import (
	"strconv"

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
	Username    string
	FirstName   *string
	LastName    *string
	Assignments []AssignmentScore
}

type AssignmentScore struct {
	ChapterID  uuid.UUID
	ItemID     int
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
	page string,
	pageSize string,
) int64 {
	defaultPage := 1
	defaultPageSize := 10

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = defaultPage
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = defaultPageSize
	}

	baseQuery := studentRepository.DB.Model(&models.Student{}).Preload("User").
		Preload("Assignments").
		Where("group_id = ?", groupId)

	var totalStudents int64
	baseQuery.Count(&totalStudents)

	offset := (pageInt - 1) * pageSizeInt
	baseQuery.Offset(offset).Limit(pageSizeInt).Find(student)

	return totalStudents
}

func (studentRepository *StudentRepository) GetStudentInGroupID(
	students *[]models.Student,
	groupID uuid.UUID,
) {
	studentRepository.DB.Where("group_id = ?", groupID).Preload("User").Find(students)
}

func (studentRepository *StudentRepository) GetStudentGroupCount(groupID uuid.UUID) int64 {
	var count int64
	studentRepository.DB.Model(models.Student{}).Where("group_id = ?", groupID).Count(&count)
	return count
}
