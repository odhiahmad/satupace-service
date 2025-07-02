package response

type ProductVariantResponse struct {
	Id          int      `json:"id"`
	SKU         string   `json:"sku,omitempty"`
	Name        string   `json:"nama"`
	Image       string   `json:"image"`
	BasePrice   *float64 `json:"base_price"`
	TrackStock  bool     `json:"track_stock"`
	Stock       int      `json:"stock,omitempty"`
	IsAvailable bool     `json:"is_available,omitempty"`
	IsActive    bool     `json:"is_active,omitempty"`
}
