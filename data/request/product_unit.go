package request

type ProductUnitCreate struct {
	Name string `json:"name" validate:"required"`
}

type ProductUnitUpdate struct {
	Id   int    `validate:"required"`
	Name string `json:"name" validate:"required"`
}
