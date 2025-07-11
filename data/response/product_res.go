package response

type ProductResponse struct {
	Id           int                      `json:"id"`
	SKU          *string                  `json:"sku,omitempty"`
	Name         string                   `json:"name"`
	Description  *string                  `json:"description,omitempty"`
	Image        *string                  `json:"image,omitempty"`
	BasePrice    *float64                 `json:"base_price,omitempty"`
	SellPrice    *float64                 `json:"sell_price,omitempty"`
	FinalPrice   *float64                 `json:"final_price,omitempty"`
	Stock        *int                     `json:"stock,omitempty"`
	TrackStock   bool                     `json:"track_stock"`
	MinimumSales *int                     `json:"minimum_sales,omitempty"`
	IsAvailable  bool                     `json:"is_available"`
	IsActive     bool                     `json:"is_active"`
	HasVariant   bool                     `json:"has_variant"`
	Variants     []ProductVariantResponse `json:"variants"`
	Category     *CategoryResponse        `json:"category"`
	Brand        *BrandResponse           `json:"brand"`
	Tax          *TaxResponse             `json:"tax"`
	Discount     *DiscountResponse        `json:"discount"`
	Unit         *UnitResponse            `json:"unit"`
	CreatedAt    string                   `json:"created_at"` // Gunakan ISO8601 format saat mapping
	UpdatedAt    string                   `json:"updated_at"`
}

type CategoryResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId *int   `json:"parent_id,omitempty"`
}
