package request

type ProductCreate struct {
	BusinessId        int                    `json:"business_id" validate:"required"`
	ProductCategoryId int                    `json:"product_category_id" validate:"required"`
	Name              string                 `json:"nama" validate:"required"`
	Description       string                 `json:"deskripsi,omitempty"`
	Image             string                 `json:"gambar" validate:"required"`
	BasePrice         float64                `json:"base_price" validate:"required"`
	FinalPrice        float64                `json:"final_price,omitempty"`
	Discount          float64                `json:"discount,omitempty"`
	Promo             float64                `json:"promo,omitempty"`
	SKU               string                 `json:"sku,omitempty"`
	Stock             int                    `json:"stock,omitempty"`
	HasVariant        bool                   `json:"has_variant"`
	Variants          []ProductVariantCreate `json:"variants,omitempty"`
}

type ProductUpdate struct {
	Id                int                    `validate:"required"`
	BusinessId        int                    `json:"business_id" validate:"required"`
	ProductCategoryId int                    `json:"product_category_id" validate:"required"`
	Name              string                 `json:"nama" validate:"required"`
	Description       string                 `json:"deskripsi,omitempty"`
	Image             string                 `json:"gambar" validate:"required"`
	BasePrice         float64                `json:"base_price" validate:"required"`
	FinalPrice        float64                `json:"final_price,omitempty"`
	Discount          float64                `json:"discount,omitempty"`
	Promo             float64                `json:"promo,omitempty"`
	SKU               string                 `json:"sku,omitempty"`
	Stock             int                    `json:"stock,omitempty"`
	IsAvailable       bool                   `json:"is_available,omitempty"`
	IsActive          bool                   `json:"is_active,omitempty"`
	HasVariant        bool                   `json:"has_variant"`
	Variants          []ProductVariantUpdate `json:"variants,omitempty"`
}
