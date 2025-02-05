package repositories

import (
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type DepartmentRepositoryQ interface{}

type DepartmentRepository struct {
	DB *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{DB: db}
}

func (deptRepository *DepartmentRepository) GetAllDepts(depts *[]models.Department) {
	deptRepository.DB.Find(depts)
}

func (deptRepository *DepartmentRepository) GetAllDeptsWithTotalMarks(depts *[]models.DepartmentWithAggregate, year string) error {
	baseQuery := deptRepository.DB.Model(models.Department{}).
		Select(`
            departments.*, 
            COALESCE(SUM(student_assignment_chapter_items.marking), 0) AS total_marks,
            COUNT(DISTINCT students.stu_id) as student_count
        `)

	joinCondition := "class_schedules.dept_id = departments.dept_id"
	if year != "" {
		joinCondition += " AND class_schedules.year = '" + year + "'"
	}

	err := baseQuery.
		Joins("LEFT JOIN class_schedules ON " + joinCondition).
		Joins("LEFT JOIN students ON students.group_id = class_schedules.group_id").
		Joins("LEFT JOIN student_assignment_chapter_items ON student_assignment_chapter_items.stu_id = students.stu_id").
		Group("departments.dept_id").
		Order("total_marks DESC").
		Find(&depts)

	if err.Error != nil {
		return err.Error
	}
	return nil
}
