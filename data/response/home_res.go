package response

import "github.com/google/uuid"

type TodaySummaryResponse struct {
	TotalRevenue float64 `json:"total_revenue"`
	TotalOrders  int64   `json:"total_orders"`
	TotalItems   int64   `json:"total_items"`
}

type HomeResponse struct {
	TodaySummary     TodaySummaryResponse      `json:"today_summary"`
	RecentOrders     []TransactionResponse     `json:"recent_orders"`
	RecentItems      []TransactionItemResponse `json:"recent_items"`
	RecentItemsTotal int64                     `json:"recent_items_total"`
	TopProducts      []TopProductResponse      `json:"top_products"`
}

type TopProductResponse struct {
	ProductId   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	OrderCount  int64     `json:"order_count"`
}
