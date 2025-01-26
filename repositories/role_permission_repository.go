package repositories

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"gorm.io/gorm"
)

type RolePermissionRepository struct {
	DB *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) *RolePermissionRepository {
	return &RolePermissionRepository{DB: db}
}

func (rolePermissionRepo *RolePermissionRepository) GetPermissionByRole(rolePermission *[]models.RolePermission, role string) {
	rolePermissionRepo.DB.Where("role = ?", role).Find(rolePermission)
}
