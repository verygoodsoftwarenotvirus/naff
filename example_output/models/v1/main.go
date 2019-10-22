package models

import "fmt"

const (
	SortAscending  sortType = "asc"
	SortDescending sortType = "desc"
)

type (
	ContextKey string
	sortType   string
	Pagination struct {
		Page       uint64
		Limit      uint64
		TotalCount uint64
	}
	CountResponse struct {
		Count uint64
	}
)

var _ error = (*ErrorResponse)(nil)

type ErrorResponse struct {
	Message string
	Code    uint
}

func (er *ErrorResponse) Error() string {
	return fmt.Sprintf("%d - %s", er.Code, er.Message)
}
