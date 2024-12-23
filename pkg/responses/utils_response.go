package responses

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Pages    int `json:"pages"`
}
