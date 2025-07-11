package request

type Pagination struct {
	Cursor     string `form:"cursor"` // encoded ID or timestamp
	Page       int    `form:"page"`   // optional fallback
	Limit      int    `form:"limit" binding:"min=1,max=100"`
	SortBy     string `form:"sort_by"`     // kolom sort (e.g. created_at)
	OrderBy    string `form:"order_by"`    // asc / desc
	Search     string `form:"search"`      // keyword pencarian
	CategoryID *int   `form:"category_id"` // opsional, filter by kategori
}
