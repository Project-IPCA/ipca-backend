package rolepermission

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"gorm.io/gorm"
)

type ServiceWrapper interface {
	Create(role, permission string) error
	Delete(rolePermission models.RolePermission) error
}

type Service struct {
	DB *gorm.DB
}

func NewRolePermissionService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
