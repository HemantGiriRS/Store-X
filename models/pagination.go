package models

type PaginationInfo struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalRecords int64 `json:"total_records"`
	TotalPages   int   `json:"total_pages"`
}
