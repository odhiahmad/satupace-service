package response

import "time"

type PromoResponse struct {
	Id               int                    `json:"id"`
	BusinessId       int                    `json:"business_id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	Type             string                 `json:"type"`
	Amount           float64                `json:"amount"`
	IsPercentage     bool                   `json:"is_percentage"`          // true = amount sebagai persen
	MinSpend         *float64               `json:"min_spend"`              // hanya untuk promo minimum spend
	MinQuantity      *int                   `json:"min_quantity"`           // hanya untuk promo minimum spend
	FreeProduct      *ProductResponse       `json:"free_product,omitempty"` // untuk Buy A+B get C, atau Buy 1 Get 1
	RequiredProducts []RequiredProductData  `json:"required_products"`
	ProductPromos    []ProductPromoResponse `json:"product_promos"`
	StartDate        time.Time              `json:"start_date"`
	EndDate          time.Time              `json:"end_date"`
	IsActive         bool                   `json:"is_active"`
}
