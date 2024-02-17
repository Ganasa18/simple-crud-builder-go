package web

type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type WebResponsePaginate struct {
	Code     int         `json:"code"`
	Status   string      `json:"status"`
	Data     interface{} `json:"list_data"`
	PageInfo interface{} `json:"page_info"`
}

type PageInfoResponse struct {
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	TotalPage   int   `json:"total_page"`
}
