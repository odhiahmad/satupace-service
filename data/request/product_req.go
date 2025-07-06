package request

type ProductRequest struct {
	BusinessId   int                     `json:"business_id" validate:"required"`
	CategoryId   *int                    `json:"category_id" validate:"required"`
	Name         string                  `json:"name" validate:"required"`
	HasVariant   bool                    `json:"has_variant"`
	Brand        *string                 `json:"brand,omitempty"`
	Description  *string                 `json:"description,omitempty"`
	Image        *string                 `json:"image,omitempty"`
	BasePrice    *float64                `json:"base_price,omitempty"`
	SellPrice    *float64                `json:"sell_price,omitempty"`
	SKU          *string                 `json:"sku,omitempty"`
	Stock        *int                    `json:"stock,omitempty"`
	TrackStock   bool                    `json:"track_stock"`
	MinimumSales *int                    `json:"minimum_sales,omitempty"`
	DiscountId   *int                    `json:"discount_id,omitempty"`
	PromoIds     []int                   `json:"promo_ids"`
	BrandId      *int                    `json:"brand_id,omitempty"`
	TaxId        *int                    `json:"tax_id,omitempty"`
	UnitId       *int                    `json:"unit_id,omitempty"`
	IsAvailable  bool                    `json:"is_available"`
	IsActive     bool                    `json:"is_active"`
	Variants     []ProductVariantRequest `json:"variants,omitempty"`
}

type CategoryRequest struct {
	Name       string `json:"name" validate:"required"`
	BusinessId int    `json:"business_id" validate:"required"`
	ParentId   *int   `json:"parent_id,omitempty"`
}

type ProductVariantRequest struct {
	BusinessId int      `json:"business_id" validate:"required"`
	Name       string   `json:"name" validate:"required"`
	BasePrice  *float64 `json:"base_price,omitempty"`
	SellPrice  *float64 `json:"sale_price,omitempty"`
	SKU        string   `json:"sku,omitempty"`
	Stock      int      `json:"stock,omitempty"`
	TrackStock *bool    `json:"track_stock,omitempty"`
}

type ProductPromoCreate struct {
	BusinessId  int `json:"business_id" validate:"required"`
	ProductId   int `json:"product_id" validate:"required"`
	PromoId     int `json:"promo_id" validate:"required"`
	MinQuantity int `json:"min_quantity"` // opsional
}

type ProductPromoUpdate struct {
	MinQuantity int `json:"min_quantity"`
}
