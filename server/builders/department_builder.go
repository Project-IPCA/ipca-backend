package builders

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type DepartmentBuilder struct {
	DeptID uuid.UUID
	Name   string
}

func NewDepartmentBuilder() *DepartmentBuilder {
	return &DepartmentBuilder{}
}

func (departmentBuilder *DepartmentBuilder) SetDeptID(deptId uuid.UUID) *DepartmentBuilder {
	departmentBuilder.DeptID = deptId
	return departmentBuilder
}

func (departmentBuilder *DepartmentBuilder) SetName(name string) *DepartmentBuilder {
	departmentBuilder.Name = name
	return departmentBuilder
}

func (departmentBuilder *DepartmentBuilder) Build() models.Department {
	department := models.Department{
		DeptID: departmentBuilder.DeptID,
		Name:   departmentBuilder.Name,
	}
	return department
}
