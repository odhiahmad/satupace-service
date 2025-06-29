package response

import "time"

type PromoResponse struct {
	Id            int                    `json:"id"`
	BusinessId    int                    `json:"business_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Type          string                 `json:"type"` // "percentage", "fixed"
	Amount        float64                `json:"amount"`
	MinQuantity   int                    `json:"min_quantity"`
	IsGlobal      bool                   `json:"is_global"`
	StartDate     time.Time              `json:"start_date"`
	EndDate       time.Time              `json:"end_date"`
	IsActive      bool                   `json:"is_active"`
	ProductPromos []ProductPromoResponse `json:"product_promos"`
}
