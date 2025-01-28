package user

import (
	"fmt"
	"github.com/Project-IPCA/ipca-backend/models"
)

func (service *Service) Delete(
	user *models.User,
) error {
	err := service.DB.Delete(*user)
	if err.Error != nil {
		return fmt.Errorf("error while delete item : %v", err.Error)
	}
	return nil
}

func (service *Service) DeleteAdmin(
	user *models.User,
) error {
	user.IsActive = false
	err := service.DB.Save(*user)
	if err.Error != nil {
		return fmt.Errorf("error while delete admin : %v", err.Error)
	}
	return nil
}
