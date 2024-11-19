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
