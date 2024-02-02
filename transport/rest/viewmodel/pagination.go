package viewmodel

type ReturnPagination struct {
	TotalRecords   int64 `json:"total_records"`
	RecordsPerPage int64 `json:"records_per_page"`
	TotalPages     int64 `json:"total_pages"`
	CurrentPage    int64 `json:"current_page"`
}

type PaginatedResult[T any] struct {
	Pagination ReturnPagination `json:"pagination,omitempty"`
	List       T                `json:"data,omitempty"`
}
