package request

type BusinessTypeCreate struct {
	Name string `json:"name" validate:"required"`
}

type BusinessTypeUpdate struct {
	Name string `json:"name" validate:"required"`
}
