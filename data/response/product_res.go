package response

type ProductResponse struct {
	Id               int                      `json:"id"`
	SKU              *string                  `json:"sku"`
	Name             string                   `json:"name"`
	Description      *string                  `json:"description"`
	Image            *string                  `json:"image"`
	BasePrice        *float64                 `json:"base_price"`
	SellPrice        *float64                 `json:"sell_price"`
	FinalPrice       *float64                 `json:"final_price"`
	Stock            *int                     `json:"stock"`
	TrackStock       bool                     `json:"track_stock"`
	IgnoreStockCheck *bool                    `json:"ignore_stock_check"`
	MinimumSales     *int                     `json:"minimum_sales"`
	IsAvailable      bool                     `json:"is_available"`
	IsActive         bool                     `json:"is_active"`
	HasVariant       bool                     `json:"has_variant"`
	Variants         []ProductVariantResponse `json:"variants"`
	Category         *CategoryResponse        `json:"category"`
	Brand            *BrandResponse           `json:"brand"`
	Tax              *TaxResponse             `json:"tax"`
	Discount         *DiscountResponse        `json:"discount"`
	Unit             *UnitResponse            `json:"unit"`
	CreatedAt        string                   `json:"created_at"` // Gunakan ISO8601 format saat mapping
	UpdatedAt        string                   `json:"updated_at"`
}

type ProductSearchResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CategoryResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId *int   `json:"parent_id"`
}
