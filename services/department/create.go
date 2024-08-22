package department

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/server/builders"
)

func (departmentService *Service) Create(name string) error {
	department := builders.NewDepartmentBuilder().
		SetDeptID(uuid.New()).
		SetName(name).
		Build()

	departmentService.DB.Create(&department)
	return nil
}
