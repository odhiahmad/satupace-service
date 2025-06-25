package request

type BundleItemRequest struct {
	ProductId int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

type BundleCreate struct {
	BusinessId  int                 `json:"business_id" validate:"required"`
	Name        string              `json:"name" validate:"required"`
	Description *string             `json:"description,omitempty"`
	Image       *string             `json:"image" validate:"required"`
	BasePrice   float64             `json:"base_price" validate:"required"`
	Discount    float64             `json:"discount,omitempty"`
	Promo       float64             `json:"promo,omitempty"`
	Items       []BundleItemRequest `json:"items" validate:"required,dive"`
}

type BundleUpdate struct {
	BusinessId  int                 `json:"business_id" validate:"required"`
	Name        string              `json:"name" validate:"required"`
	Description *string             `json:"description,omitempty"`
	Image       *string             `json:"image" validate:"required"`
	BasePrice   float64             `json:"base_price" validate:"required"`
	Items       []BundleItemRequest `json:"items" validate:"required,dive"`
	IsAvailable bool                `json:"is_available,omitempty"`
	IsActive    bool                `json:"is_active,omitempty"`
}
