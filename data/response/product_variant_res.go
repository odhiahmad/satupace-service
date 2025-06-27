package response

type ProductVariantResponse struct {
	Id          int      `json:"id"`
	Name        string   `json:"nama"`
	Image       string   `json:"gambar"`
	BasePrice   *float64 `json:"base_price"`
	FinalPrice  *float64 `json:"final_price,omitempty"`
	TaxId       *int     `json:"tax_id,omitempty"`
	DiscountId  *int     `json:"discount_id,omitempty"`
	PromoIds    []int    `json:"promo_ids"`
	SKU         *string  `json:"sku,omitempty"`
	Stock       int      `json:"stock,omitempty"`
	IsAvailable bool     `json:"is_available,omitempty"`
	IsActive    bool     `json:"is_active,omitempty"`
}
