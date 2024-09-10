package classlabstaff

import "github.com/Project-IPCA/ipca-backend/models"

func (classLabStaffService *Service) DeleteAll(classLabStaff *[]models.ClassLabStaff) {
	classLabStaffService.DB.Delete(classLabStaff)
}
