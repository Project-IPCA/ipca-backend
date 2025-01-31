package rolepermission

import (
	"fmt"

	"github.com/Project-IPCA/ipca-backend/models"
)

func (service *Service) Delete(
	rolePermission *models.RolePermission,
) error {
	err := service.DB.Delete(*rolePermission)
	if err.Error != nil {
		return fmt.Errorf("error while delete role permission : %v", err.Error)
	}
	return nil
}
