package responses

import (
	"github.com/google/uuid"

	"github.com/Project-IPCA/ipca-backend/models"
)

type DepartmentResponse struct {
	DeptID uuid.UUID `json:"dept_id"`
	NameTH string    `json:"name_th"`
	NameEN string    `json:"name_en"`
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

type AverageDeptScoreResponse struct {
	MaxRange int                    `json:"max_range"`
	Data     []AverageDeptScoreData `json:"data"`
}

type AverageDeptScoreData struct {
	NameTH string  `json:"dept_name_th"`
	NameEN string  `json:"dept_name_en"`
	Score  float64 `json:"score"`
}

func NewAverageDeptScoreResponse(depts []models.DepartmentWithAggregate,labClassInfo []models.LabClassInfo) AverageDeptScoreResponse {
	data := make([]AverageDeptScoreData, 0)
	for _, dept := range depts {
		var average float64
		if dept.StudentCount == 0 {
			average = float64(dept.TotalMarks)
		} else {
			average = float64(dept.TotalMarks) / float64(dept.StudentCount)
		}
		data = append(data, AverageDeptScoreData{
			NameTH: dept.NameTH,
			NameEN: dept.NameEN,
			Score:  average,
		})
	}
	
	response := AverageDeptScoreResponse{
		MaxRange: len(labClassInfo)*labClassInfo[0].FullMark,
		Data: data,
	}

	return response
}
