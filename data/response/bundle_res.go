package response

type BundleItemResponse struct {
	Id        int    `json:"id"`
	ProductId int    `json:"product_id"`
	Product   string `json:"product"` // bisa diganti sesuai kebutuhan
	Quantity  int    `json:"quantity"`
}

type BundleResponse struct {
	Id                int                      `json:"id"`
	Name              string                   `json:"name"`
	Description       string                   `json:"description,omitempty"`
	Image             string                   `json:"image"`
	BasePrice         float64                  `json:"base_price"`
	FinalPrice        float64                  `json:"final_price,omitempty"`
	Discount          float64                  `json:"discount,omitempty"`
	Promo             float64                  `json:"promo,omitempty"`
	SKU               string                   `json:"sku,omitempty"`
	Stock             int                      `json:"stock,omitempty"`
	IsAvailable       bool                     `json:"is_available"`
	IsActive          bool                     `json:"is_active"`
	HasVariant        bool                     `json:"has_variant,omitempty"`
	Variants          []ProductVariantResponse `json:"variants,omitempty"`
	Items             []BundleItemResponse     `json:"items"`
	ProductCategoryId int                      `json:"product_category_id,omitempty"`
	ProductCategory   *ProductCategoryResponse `json:"product_category,omitempty"`
}
