package request

type BundleItemRequest struct {
	ProductId int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,gte=1"`
}

type BundleRequest struct {
	BusinessId  int                 `json:"business_id" validate:"required"`
	Name        string              `json:"name" validate:"required"`
	Description *string             `json:"description,omitempty"`
	Image       *string             `json:"image" validate:"required"`
	BasePrice   *float64            `json:"base_price,omitempty"`
	SellPrice   *float64            `json:"sell_price,omitempty"`
	Stock       int                 `json:"stock,omitempty"`
	TaxId       *int                `json:"tax_id,omitempty"`
	Items       []BundleItemRequest `json:"items" validate:"required,dive"`
}
