package builders

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type DepartmentBuilder struct {
	DeptID uuid.UUID
	NameTH string
	NameEN string
}

func NewDepartmentBuilder() *DepartmentBuilder {
	return &DepartmentBuilder{}
}

func (departmentBuilder *DepartmentBuilder) SetDeptID(deptId uuid.UUID) *DepartmentBuilder {
	departmentBuilder.DeptID = deptId
	return departmentBuilder
}

func (departmentBuilder *DepartmentBuilder) SetName(name_th string) *DepartmentBuilder {
	departmentBuilder.NameTH = name_th
	return departmentBuilder
}

func (departmentBuilder *DepartmentBuilder) SetNameEN(name_en string) *DepartmentBuilder {
	departmentBuilder.NameEN = name_en
	return departmentBuilder
}

func (departmentBuilder *DepartmentBuilder) Build() models.Department {
	department := models.Department{
		DeptID: departmentBuilder.DeptID,
		NameTH: departmentBuilder.NameTH,
		NameEN: departmentBuilder.NameEN,
	}
	return department
}
