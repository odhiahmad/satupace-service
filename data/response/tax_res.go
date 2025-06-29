package response

type TaxResponse struct {
	Id         int               `json:"id"`
	BusinessId int               `json:"business_id"`
	Name       string            `json:"name"`
	Type       string            `json:"type"` // "percentage", "fixed"
	Amount     float64           `json:"amount"`
	IsGlobal   bool              `json:"is_global"`
	IsActive   bool              `json:"is_active"`
	Products   []ProductResponse `json:"products,omitempty"` // hanya jika IsGlobal == false
}
