package builders

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/google/uuid"
)

type RolePermissionBuilder struct {
	PermissionID uuid.UUID
	Role         string
	Permission   string
}

func NewRolePermissionBuilder() *RolePermissionBuilder {
	return &RolePermissionBuilder{}
}

func (rolePermissionBuilder *RolePermissionBuilder) SetPermissionID(permisisonId uuid.UUID) *RolePermissionBuilder {
	rolePermissionBuilder.PermissionID = permisisonId
	return rolePermissionBuilder
}

func (rolePermissionBuilder *RolePermissionBuilder) SetRole(role string) *RolePermissionBuilder {
	rolePermissionBuilder.Role = role
	return rolePermissionBuilder
}

func (rolePermissionBuilder *RolePermissionBuilder) SetPermission(permisison string) *RolePermissionBuilder {
	rolePermissionBuilder.Permission = permisison
	return rolePermissionBuilder
}

func (rolePermissionBuilder *RolePermissionBuilder) Build() models.RolePermission {
	rolePermission := models.RolePermission{
		PermissionID: rolePermissionBuilder.PermissionID,
		Role:         rolePermissionBuilder.Role,
		Permission:   rolePermissionBuilder.Permission,
	}
	return rolePermission
}
