package responses

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type DepartmentResponse struct {
	DeptID uuid.UUID `json:"dept_id"`
	Name   string    `json:"dept_name"`
}

func NewDepartmentsResponse(departments []models.Department) *[]DepartmentResponse {
	departmentsResponse := make([]DepartmentResponse, 0)
	for _, dept := range departments {
		departmentsResponse = append(departmentsResponse, DepartmentResponse{
			DeptID: dept.DeptID,
			Name:   dept.Name,
		})
	}
	return &departmentsResponse
}
