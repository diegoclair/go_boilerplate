package viewmodel

import "math"

type ReturnPagination struct {
	TotalRecords   int64 `json:"total_records"`
	RecordsPerPage int64 `json:"records_per_page"`
	TotalPages     int64 `json:"total_pages"`
	CurrentPage    int64 `json:"current_page"`
}

type PaginatedResponse[T any] struct {
	Pagination ReturnPagination `json:"pagination"`
	List       T                `json:"data"`
}

// BuildPaginatedResponse builds a paginated result based on the given parameters.
// It takes a list of type T, the number of records to skip, the number of records to take,
// and the total number of records available.
// It returns a PaginatedResult of type T, which contains the paginated list and pagination information.
func BuildPaginatedResponse[T any](list T, skip int64, take int64, totalRecords int64) PaginatedResponse[T] {
	return PaginatedResponse[T]{
		List: list,
		Pagination: ReturnPagination{
			CurrentPage:    (skip / take) + 1,
			RecordsPerPage: take,
			TotalRecords:   totalRecords,
			TotalPages:     int64(math.Ceil(float64(totalRecords) / float64(take))),
		},
	}
}
