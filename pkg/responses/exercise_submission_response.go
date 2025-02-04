package responses

type TotalSubmissionsResponse struct {
	TotalSubmission int64 `json:"total_submission"`
}

func NewTotalSubmissionsResponse(total int64) TotalSubmissionsResponse {
	response := TotalSubmissionsResponse{
		TotalSubmission: total,
	}
	return response
}
