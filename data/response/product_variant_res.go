package response

type ProductVariantResponse struct {
	Id               int      `json:"id"`
	SKU              *string  `json:"sku,omitempty"`
	Name             string   `json:"name"`
	Description      *string  `json:"description,omitempty"`
	Image            string   `json:"image"`
	BasePrice        *float64 `json:"base_price,omitempty"`
	SellPrice        *float64 `json:"sell_price,omitempty"`
	FinalPrice       *float64 `json:"final_price,omitempty"`
	TrackStock       bool     `json:"track_stock"`
	IgnoreStockCheck *bool    `json:"ignore_stock_check"`
	Stock            int      `json:"stock,omitempty"`
	IsAvailable      bool     `json:"is_available,omitempty"`
	IsActive         bool     `json:"is_active,omitempty"`
}
