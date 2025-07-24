package response

type TaxResponse struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	IsPercentage *bool   `json:"is_percentage"`
	Amount       float64 `json:"amount"`
	IsGlobal     *bool   `json:"is_global"`
	IsActive     *bool   `json:"is_active"`
}
