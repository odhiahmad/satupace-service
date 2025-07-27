package response

type BundleItemResponse struct {
	Id          int      `json:"id"`
	ProductId   int      `json:"product_id"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Image       *string  `json:"image"`
	BasePrice   *float64 `json:"base_price"`
	SellPrice   *float64 `json:"sell_price"`
	SKU         *string  `json:"sku"`
	Stock       *int     `json:"stock"`
	IsAvailable bool     `json:"is_available"`
	IsActive    bool     `json:"is_active"`
	Quantity    int      `json:"quantity"`
}

type BundleResponse struct {
	Id          int                  `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Image       string               `json:"image"`
	BasePrice   *float64             `json:"base_price"`
	SellPrice   *float64             `json:"sell_price"`
	SKU         string               `json:"sku"`
	Stock       int                  `json:"stock"`
	IsAvailable bool                 `json:"is_available"`
	IsActive    bool                 `json:"is_active"`
	Items       []BundleItemResponse `json:"items"`
}
