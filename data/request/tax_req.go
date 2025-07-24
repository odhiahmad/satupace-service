package request

type TaxRequest struct {
	BusinessId int     `json:"business_id" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
	IsGlobal   *bool   `json:"is_global"`
	IsActive   *bool   `json:"is_active"`
}
