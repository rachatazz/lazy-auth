package model

type QueryPagination struct {
	OrderBy *string `form:"order_by"`
	SortBy  *string `form:"sort_by"`
	Limit   *int    `form:"limit"`
	Offset  *int    `form:"offset"`
}

type MetaPagination struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
