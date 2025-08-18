package request

type OrderTypeRequest struct {
	Name string `json:"name" validate:"required"`
}
