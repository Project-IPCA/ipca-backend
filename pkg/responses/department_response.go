package responses

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type DepartmentResponse struct {
	DeptID uuid.UUID `json:"dept_id"`
	NameTH string    `json:"dept_name"`
	NameEN string    `json:"dept_name_en"`
}

func NewDepartmentsResponse(departments []models.Department) *[]DepartmentResponse {
	departmentsResponse := make([]DepartmentResponse, 0)
	for _, dept := range departments {
		departmentsResponse = append(departmentsResponse, DepartmentResponse{
			DeptID: dept.DeptID,
			NameTH: dept.NameTH,
			NameEN: dept.NameEN,
		})
	}
	return &departmentsResponse
}
