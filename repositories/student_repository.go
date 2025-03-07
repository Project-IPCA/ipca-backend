package repositories

import (
	"strconv"
	"strings"

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
	studentRepository.DB.Where("stu_id = ?", stuId).Preload("Group").Find(student)
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

func (studentRepository *StudentRepository) GetStudentGroupOrYearCount(groupID uuid.UUID, year string, language string) int64 {
	var count int64
	baseQuery := studentRepository.DB.Model(models.Student{}).Joins("JOIN class_schedules ON class_schedules.group_id = students.group_id")

	if groupID != uuid.Nil {
		baseQuery.Where("students.group_id = ?", groupID)
	}

	if language != "" {
		baseQuery.Where("class_schedules.language = ?", strings.ToUpper(language))
	}

	if year != "" {
		baseQuery.Where("class_schedules.year = ?", year)
	}

	baseQuery.Count(&count)
	return count
}

func (studentRepository *StudentRepository) GetTotalStudent(groupId string, year string, status string, language string) int64 {
	var total int64
	baseQuery := studentRepository.DB.Model(models.Student{}).Joins("JOIN class_schedules ON class_schedules.group_id = students.group_id")
	if status != "" {
		baseQuery = baseQuery.Joins("JOIN users ON users.user_id = students.stu_id").Where("users.is_online = ?", status)
	}

	if groupId != "" {
		baseQuery = baseQuery.Where("students.group_id = ? ", groupId)
	}

	if year != "" {
		baseQuery = baseQuery.Where("class_schedules.year = ?", year)
	}

	if language != "" {
		baseQuery = baseQuery.Where("class_schedules.language = ?", strings.ToUpper(language))
	}
	baseQuery.Count(&total)
	return total
}

func (studentRepository *StudentRepository) GetStudentGroupRanking(students *[]models.StudentWithAggregate, groupId uuid.UUID) error {
	err := studentRepository.DB.Model(&models.Student{}).
		Preload("User").
		Joins("LEFT JOIN student_assignment_chapter_items ON student_assignment_chapter_items.stu_id = students.stu_id").
		Where("students.group_id = ?", groupId).
		Select("students.*, COALESCE(SUM(student_assignment_chapter_items.marking), 0) AS total_marks").
		Group("students.stu_id").
		Order("total_marks DESC").
		Find(students)

	if err.Error != nil {
		return err.Error
	}
	return nil
}
