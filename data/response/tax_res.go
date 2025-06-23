package response

import "time"

type TaxResponse struct {
	Id          int               `json:"id"`
	BusinessId  int               `json:"business_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"type"` // "percentage", "fixed"
	Amount      float64           `json:"amount"`
	IsGlobal    bool              `json:"is_global"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Products    []ProductResponse `json:"products,omitempty"` // hanya jika IsGlobal == false
}
