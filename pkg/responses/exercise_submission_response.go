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

type StatsSubmissionsResponse struct {
	SubmissionsList []int64  `json:"submissions_list"`
	DateList        []string `json:"date_list"`
}

func NewStatsSubmissionsResponse(
	submissionList []int64,
	dateList []string,
) StatsSubmissionsResponse {
	return StatsSubmissionsResponse{
		SubmissionsList: submissionList,
		DateList:        dateList,
	}
}
