package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Project-IPCA/ipca-backend/models"
)

type ClassLabStaffRepositoryQ interface{}

type ClassLabStaffRepository struct {
	DB *gorm.DB
}

func NewClassLabStaffRepository(db *gorm.DB) *ClassLabStaffRepository {
	return &ClassLabStaffRepository{DB: db}
}

func (repo *ClassLabStaffRepository) GetClassLabStaffByGroupID(
	classLabStaff *[]models.ClassLabStaff,
	groupId uuid.UUID,
) {
	repo.DB.Where("class_id = ?", groupId).Find(classLabStaff)
}

func (repo *ClassLabStaffRepository) CheckStaffValidInClass(
	groupId uuid.UUID,
	staffId uuid.UUID,
) bool {
	var classLabStaff models.ClassLabStaff
	err := repo.DB.Where("class_id = ? AND staff_id = ?", groupId, staffId).First(classLabStaff)
	return err.Error == nil
}
