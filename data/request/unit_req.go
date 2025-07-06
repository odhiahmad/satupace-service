package request

type UnitRequest struct {
	BusinessId int     `json:"business_id" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	Alias      string  `json:"alias"`
	Multiplier float64 `json:"multiplier" validate:"required,gte=1"`
}
