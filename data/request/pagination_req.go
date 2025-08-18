package request

type Pagination struct {
	Cursor     string `form:"cursor"`
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
	SortBy     string `form:"sort_by"`
	OrderBy    string `form:"order_by"`
	Search     string `form:"search"`
	CategoryID *int   `form:"category_id"`
}
