package request

type PaymentMethodRequest struct {
	Name string `json:"name" validate:"required"`
}
