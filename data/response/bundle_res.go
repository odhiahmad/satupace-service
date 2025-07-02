package response

type BundleItemResponse struct {
	Id          int     `json:"id"`
	ProductId   int     `json:"product_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Image       *string `json:"image,omitempty"`
	BasePrice   float64 `json:"base_price"`
	SKU         *string `json:"sku,omitempty"`
	Stock       *int    `json:"stock,omitempty"`
	IsAvailable bool    `json:"is_available,omitempty"`
	IsActive    bool    `json:"is_active,omitempty"`
	Quantity    int     `json:"quantity"`
}

type BundleResponse struct {
	Id                int                      `json:"id"`
	Name              string                   `json:"name"`
	Description       string                   `json:"description,omitempty"`
	Image             string                   `json:"image"`
	BasePrice         float64                  `json:"base_price"`
	SKU               string                   `json:"sku,omitempty"`
	Stock             int                      `json:"stock,omitempty"`
	IsAvailable       bool                     `json:"is_available"`
	IsActive          bool                     `json:"is_active"`
	Items             []BundleItemResponse     `json:"items"`
	ProductCategoryId int                      `json:"product_category_id,omitempty"`
	ProductCategory   *ProductCategoryResponse `json:"product_category,omitempty"`
}
