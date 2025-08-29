package request

type OrderTypeRequest struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}
