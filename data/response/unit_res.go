package response

type UnitResponse struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"` // "Pcs", "Kg", dll
	Alias      string  `json:"alias"`
	Multiplier float64 `json:"multiplier" validate:"required,gte=1"`
}
