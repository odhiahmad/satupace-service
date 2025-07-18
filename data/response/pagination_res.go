package response

type PaginatedResponse struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Total     int64  `json:"total"`
	OrderBy   string `json:"order_by"`
	SortOrder string `json:"sort_by"`
}

type CursorPaginatedResponse struct {
	Limit      int    `json:"limit"`
	SortBy     string `json:"sort_by"`
	OrderBy    string `json:"order_by"`
	NextCursor string `json:"next_cursor,omitempty"`
}
