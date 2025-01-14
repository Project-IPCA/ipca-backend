package classschedule

import (
	"fmt"
	"github.com/Project-IPCA/ipca-backend/models"
)

func (service *Service) Delete(
	classSchedule *models.ClassSchedule,
) error {
	err := service.DB.Delete(*classSchedule)
	if err.Error != nil {
		return fmt.Errorf("error while delete item : %v", err.Error)
	}
	return nil
}
