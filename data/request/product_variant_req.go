package request

type ProductVariantCreate struct {
	BusinessId int     `json:"business_id" validate:"required"`
	Name       string  `json:"nama" validate:"required"`
	Image      string  `json:"gambar" validate:"required"`
	BasePrice  float64 `json:"base_price" validate:"required"`
	FinalPrice float64 `json:"final_price,omitempty"`
	Discount   float64 `json:"discount,omitempty"`
	Promo      float64 `json:"promo,omitempty"`
	SKU        string  `json:"sku,omitempty"`
	Stock      int     `json:"stock,omitempty"`
}

type ProductVariantUpdate struct {
	Id          int     `validate:"required"`
	BusinessId  int     `json:"business_id" validate:"required"`
	Name        string  `json:"nama" validate:"required"`
	Image       string  `json:"gambar" validate:"required"`
	BasePrice   float64 `json:"base_price" validate:"required"`
	FinalPrice  float64 `json:"final_price,omitempty"`
	Discount    float64 `json:"discount,omitempty"`
	Promo       float64 `json:"promo,omitempty"`
	SKU         string  `json:"sku,omitempty"`
	Stock       int     `json:"stock,omitempty"`
	IsAvailable bool    `json:"is_available,omitempty"`
	IsActive    bool    `json:"is_active,omitempty"`
}
