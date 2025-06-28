package request

type ProductCreate struct {
	BusinessId        int                    `json:"business_id" validate:"required"`
	ProductCategoryId *int                   `json:"product_category_id" validate:"required"`
	Name              string                 `json:"name" validate:"required"`
	Description       *string                `json:"description,omitempty"`
	Image             *string                `json:"image" validate:"required"`
	BasePrice         *float64               `json:"base_price" validate:"required,gt=0"`
	FinalPrice        *float64               `json:"final_price,omitempty"`
	SKU               *string                `json:"sku,omitempty"`
	Stock             *int                   `json:"stock,omitempty"`
	TrackStock        *bool                  `json:"track_stock,omitempty"`
	Variants          []ProductVariantCreate `json:"variants,omitempty"`
	TaxId             *int                   `json:"tax_id,omitempty"`
	DiscountId        *int                   `json:"discount_id,omitempty"`
	PromoIds          []int                  `json:"promo_ids"`
	UnitId            *int                   `json:"unit_id,omitempty"`
}

type ProductUpdate struct {
	ProductCategoryId *int                   `json:"product_category_id" validate:"required"`
	Name              string                 `json:"name" validate:"required"`
	Description       *string                `json:"description,omitempty"`
	Image             *string                `json:"image" validate:"required"`
	BasePrice         *float64               `json:"base_price" validate:"required,gt=0"`
	FinalPrice        *float64               `json:"final_price,omitempty"`
	SKU               *string                `json:"sku,omitempty"`
	Stock             *int                   `json:"stock,omitempty"`
	TrackStock        *bool                  `json:"track_stock,omitempty"`
	Variants          []ProductVariantUpdate `json:"variants,omitempty"`
	TaxId             *int                   `json:"tax_id,omitempty"`
	DiscountId        *int                   `json:"discount_id,omitempty"`
	PromoIds          []int                  `json:"promo_ids"`
	UnitId            *int                   `json:"unit_id,omitempty"`
}

type ProductCategoryCreate struct {
	Name       string `json:"name" validate:"required"`
	BusinessId int    `json:"business_id" validate:"required"`
	ParentId   *int   `json:"parent_id,omitempty"`
}

type ProductCategoryUpdate struct {
	Id       int    `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	ParentId *int   `json:"parent_id,omitempty"`
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

type ProductUnitCreate struct {
	BusinessId int     `json:"business_id" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	Alias      string  `json:"alias"`
	Multiplier float64 `json:"multiplier" validate:"required,gte=1"`
}

type ProductUnitUpdate struct {
	Id         int     `json:"id" validate:"required"`
	BusinessId int     `json:"business_id" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	Alias      string  `json:"alias"`
	Multiplier float64 `json:"multiplier" validate:"required,gte=1"`
}

type ProductVariantCreate struct {
	BusinessId int      `json:"business_id" validate:"required"`
	Name       string   `json:"name" validate:"required"`
	Image      *string  `json:"image" validate:"required"`
	BasePrice  *float64 `json:"base_price" validate:"required"`
	SKU        string   `json:"sku,omitempty"`
	Stock      int      `json:"stock,omitempty"`
	TrackStock *bool    `json:"track_stock,omitempty"`
	TaxId      *int     `json:"tax_id,omitempty"`
	DiscountId *int     `json:"discount_id,omitempty"`
	PromoIds   []int    `json:"promo_ids"`
	UnitId     *int     `json:"unit_id,omitempty"`
}

type ProductVariantUpdate struct {
	BusinessId int      `json:"business_id" validate:"required"`
	Name       string   `json:"name" validate:"required"`
	Image      *string  `json:"image" validate:"required"`
	BasePrice  *float64 `json:"base_price" validate:"required"`
	SKU        string   `json:"sku,omitempty"`
	Stock      int      `json:"stock,omitempty"`
	TrackStock *bool    `json:"track_stock,omitempty"`
	TaxId      *int     `json:"tax_id,omitempty"`
	DiscountId *int     `json:"discount_id,omitempty"`
	PromoIds   []int    `json:"promo_ids"`
	UnitId     *int     `json:"unit_id,omitempty"`
}
