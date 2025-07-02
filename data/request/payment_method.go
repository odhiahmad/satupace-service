package request

type PaymentMethodCreate struct {
	Name string `json:"name" validate:"required"`
}

type PaymentMethodUpdate struct {
	Name string `json:"name" validate:"required"`
}
