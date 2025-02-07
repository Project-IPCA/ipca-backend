package responses

import "github.com/Project-IPCA/ipca-backend/models"

type LogLastTimeResponse struct {
	Logs  []models.ActivityLog `json:"logs"`
	Total int64                `json:"total"`
}

func NewLogLastTimeResponse(logs []models.ActivityLog, total int64) LogLastTimeResponse {
	sortLog := make([]models.ActivityLog,0)
	for i := len(logs)-1; i>=0 ; i-- {
		sortLog = append(sortLog, logs[i])
	}
	response := LogLastTimeResponse{
		Logs:  sortLog,
		Total: total,
	}
	return response
}
