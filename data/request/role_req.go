package request

type RoleCreate struct {
	Nama string `json:"nama" validate:"required"`
}

type RoleUpdate struct {
	Id   int    `validate:"required"`
	Nama string `json:"nama" validate:"required"`
}
