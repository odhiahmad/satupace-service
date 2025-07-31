package request

import "github.com/google/uuid"

type BundleItemRequest struct {
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,gte=1"`
}

type BundleRequest struct {
	BusinessId  uuid.UUID           `json:"business_id" validate:"required"`
	Name        string              `json:"name" validate:"required"`
	Description *string             `json:"description"`
	Image       *string             `json:"image" validate:"required"`
	BasePrice   *float64            `json:"base_price"`
	SellPrice   *float64            `json:"sell_price"`
	Stock       *int                `json:"stock"`
	TaxId       *uuid.UUID          `json:"tax_id"`
	Items       []BundleItemRequest `json:"items" validate:"required,dive"`
}
