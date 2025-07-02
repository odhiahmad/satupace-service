package response

import "time"

type ProductPromoResponse struct {
	IsGlobal  bool      `json:"promo_is_global"`
	StartDate time.Time `json:"promo_start_date"`
	EndDate   time.Time `json:"promo_end_date"`
}
