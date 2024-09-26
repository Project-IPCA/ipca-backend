package groupchapterselecteditem

import (
	"fmt"

	"github.com/Project-IPCA/ipca-backend/models"
)

func (service *Service) Delete (
	groupChapterSelectedItem *models.GroupChapterSelectedItem,
) error {
	err := service.DB.Delete(groupChapterSelectedItem)
	if(err!= nil){
		return fmt.Errorf("error while delete item : %v" ,err)
	}
	return nil
}