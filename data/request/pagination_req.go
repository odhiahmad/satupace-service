package request

type Pagination struct {
	Page    int    `form:"page" binding:"omitempty,min=1"`
	Limit   int    `form:"limit" binding:"omitempty,min=1,max=100"`
	SortBy  string `form:"sort_by" binding:"omitempty"`
	OrderBy string `form:"order_by" binding:"omitempty,oneof=asc desc"`
	Search  string `form:"search" binding:"omitempty"`
}
