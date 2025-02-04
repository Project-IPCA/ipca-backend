package responses

type TotalSubmissionsResponse struct {
	TotalSubmissions int64 `json:"total_submissions"`
}

func NewTotalSubmissionsResponse(total int64) TotalSubmissionsResponse {
	response := TotalSubmissionsResponse{
		TotalSubmissions: total,
	}
	return response
}
