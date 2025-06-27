package response

type ProductResponse struct {
	Id              int                      `json:"id"`
	Name            string                   `json:"nama"`
	Description     *string                  `json:"deskripsi,omitempty"`
	Image           *string                  `json:"gambar"`
	BasePrice       float64                  `json:"base_price"`
	FinalPrice      *float64                 `json:"final_price,omitempty"`
	SKU             string                   `json:"sku,omitempty"`
	Stock           int                      `json:"stock,omitempty"`
	IsAvailable     bool                     `json:"is_available,omitempty"`
	IsActive        bool                     `json:"is_active,omitempty"`
	HasVariant      bool                     `json:"has_variant"`
	ProductCategory *ProductCategoryResponse `json:"category,omitempty"`
	Variants        []ProductVariantResponse `json:"variants,omitempty"`
	Tax             *TaxResponse             `json:"tax,omitempty"`
	Discount        *DiscountResponse        `json:"discount,omitempty"`
	Promos          []ProductPromoResponse   `json:"promos"`
	Unit            *ProductUnitResponse     `json:"unit,omitempty"`
}

type ProductCategoryResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId *int   `json:"parent_id,omitempty"`
}

type ProductPromoResponse struct {
	ProductId   int    `json:"product_id"`
	ProductName string `json:"product_name"` // optional, untuk tampilkan nama
	PromoId     int    `json:"promo_id"`
	MinQuantity int    `json:"min_quantity"`
}

type ProductUnitResponse struct {
	Id         int     `json:"id"`
	BusinessId int     `json:"business_id"`
	Name       string  `json:"name"` // "Pcs", "Kg", dll
	Alias      string  `json:"alias"`
	Multiplier float64 `json:"multiplier" validate:"required,gte=1"`
}
