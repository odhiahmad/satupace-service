package response

type ProductResponse struct {
	Id                int                      `json:"id"`
	Name              string                   `json:"nama"`
	Description       string                   `json:"deskripsi,omitempty"`
	Image             string                   `json:"gambar"`
	BasePrice         float64                  `json:"base_price"`
	FinalPrice        float64                  `json:"final_price,omitempty"`
	Discount          float64                  `json:"discount,omitempty"`
	Promo             float64                  `json:"promo,omitempty"`
	SKU               string                   `json:"sku,omitempty"`
	Stock             int                      `json:"stock,omitempty"`
	IsAvailable       bool                     `json:"is_available,omitempty"`
	IsActive          bool                     `json:"is_active,omitempty"`
	HasVariant        bool                     `json:"has_variant"`
	Variants          []ProductVariantResponse `json:"variants,omitempty"`
	ProductCategory   *ProductCategoryResponse `json:"category,omitempty"`
	ProductCategoryId int                      `json:"product_category_id"`
}
