package response

import "github.com/google/uuid"

type BundleItemResponse struct {
	Id          uuid.UUID `json:"id"`
	ProductId   uuid.UUID `json:"product_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Image       *string   `json:"image"`
	BasePrice   *float64  `json:"base_price"`
	SellPrice   *float64  `json:"sell_price"`
	SKU         *string   `json:"sku"`
	Stock       *int      `json:"stock"`
	IsAvailable bool      `json:"is_available"`
	IsActive    bool      `json:"is_active"`
	Quantity    int       `json:"quantity"`
}

type BundleResponse struct {
	Id          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Image       string               `json:"image"`
	BasePrice   *float64             `json:"base_price"`
	SellPrice   *float64             `json:"sell_price"`
	SKU         string               `json:"sku"`
	Stock       *int                 `json:"stock"`
	IsAvailable bool                 `json:"is_available"`
	IsActive    bool                 `json:"is_active"`
	Items       []BundleItemResponse `json:"items"`
}
