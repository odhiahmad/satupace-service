package request

type BusinessTypeCreate struct {
	Name string `json:"name" validate:"required"`
}

type BusinessTypeUpdate struct {
	Id   int    `validate:"required"`
	Name string `json:"name" validate:"required"`
}
