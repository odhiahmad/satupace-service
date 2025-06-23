package request

type Pagination struct {
	Limit   int    `form:"limit" binding:"gte=1,lte=100"`
	Search  string `form:"search"`
	Cursor  string `form:"cursor"`   // optional: untuk cursor pagination
	SortBy  string `form:"sort_by"`  // e.g. "created_at"
	OrderBy string `form:"order_by"` // "asc" or "desc"
	Before  string `form:"before"`   // Cursor value
	After   string `form:"after"`    // Cursor value
}
