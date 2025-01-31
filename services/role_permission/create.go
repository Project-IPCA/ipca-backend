package rolepermission

import (
	"fmt"

	"github.com/Project-IPCA/ipca-backend/server/builders"
	"github.com/google/uuid"
)

func (rolePermissionService *Service) Create(
	role, permission string,
) error {
	rolePermission := builders.NewRolePermissionBuilder().
		SetPermissionID(uuid.New()).
		SetRole(role).
		SetPermission(permission).
		Build()

	err := rolePermissionService.DB.Create(&rolePermission)
	if err.Error != nil {
		return fmt.Errorf("error wihile create role permission : %v", err.Error)
	}
	return nil
}
