package response

type ProductVariantResponse struct {
	Id               int      `json:"id"`
	SKU              *string  `json:"sku"`
	Name             string   `json:"name"`
	Description      *string  `json:"description"`
	Image            string   `json:"image"`
	BasePrice        *float64 `json:"base_price"`
	SellPrice        *float64 `json:"sell_price"`
	FinalPrice       *float64 `json:"final_price"`
	TrackStock       bool     `json:"track_stock"`
	IgnoreStockCheck *bool    `json:"ignore_stock_check"`
	Stock            int      `json:"stock"`
	IsAvailable      bool     `json:"is_available"`
	IsActive         bool     `json:"is_active"`
}
