package request

type RoleCreate struct {
	Name string `json:"name" validate:"required"`
}

type RoleUpdate struct {
	Id   int    `validate:"required"`
	Name string `json:"name" validate:"required"`
}
