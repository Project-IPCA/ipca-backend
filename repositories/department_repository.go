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
