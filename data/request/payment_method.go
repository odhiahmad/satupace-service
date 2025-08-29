package request

type PaymentMethodRequest struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}
