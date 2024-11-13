package request

type PaymentMethodCreate struct {
	Name string `json:"name" validate:"required"`
}

type PaymentMethodUpdate struct {
	Id   int    `validate:"required"`
	Name string `json:"name" validate:"required"`
}
