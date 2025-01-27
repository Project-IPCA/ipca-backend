package utils

import (
	"github.com/Project-IPCA/ipca-backend/models"
	"github.com/Project-IPCA/ipca-backend/pkg/constants"
)

func ValidateRolePermission(rolePermission []models.RolePermission, checkPermission string) bool {
	for _, permission := range rolePermission {
		if permission.Permission == checkPermission {
			return true
		}
	}
	return false
}

func ValidateSupervisorAndBeyonder(user models.User) bool {
	if *user.Role == constants.Role.Supervisor || *user.Role == constants.Role.Beyonder {
		return true
	}
	return false
}

func ValidateAdminRole(user models.User) bool {
	if *user.Role == constants.Role.Supervisor || *user.Role == constants.Role.Beyonder || *user.Role == constants.Role.Ta || *user.Role == constants.Role.Executive {
		return true
	}
	return false
}
