package inout

type BaseResponse struct {
	ErrorCode        int    `json:"error_code"`
	ErrorDescription string `json:"error_description"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}