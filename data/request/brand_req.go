package request

type BrandRequest struct {
	BusinessId int    `json:"business_id" validate:"required"`
	Name       string `json:"name" validate:"required"`
}
