package request

import "github.com/google/uuid"

type ProductRequest struct {
	BusinessId       uuid.UUID               `json:"business_id" validate:"required"`
	CategoryId       *uuid.UUID              `json:"category_id"`
	Name             string                  `json:"name" validate:"required"`
	HasVariant       bool                    `json:"has_variant"`
	Brand            *string                 `json:"brand"`
	Description      *string                 `json:"description"`
	Image            *string                 `json:"image"`
	BasePrice        *float64                `json:"base_price"`
	SellPrice        *float64                `json:"sell_price"`
	SKU              *string                 `json:"sku"`
	Stock            *int                    `json:"stock"`
	TrackStock       *bool                   `json:"track_stock"`
	IgnoreStockCheck *bool                   `json:"ignore_stock_check"`
	MinimumSales     *int                    `json:"minimum_sales"`
	DiscountId       *uuid.UUID              `json:"discount_id"`
	BrandId          *uuid.UUID              `json:"brand_id"`
	TaxId            *uuid.UUID              `json:"tax_id"`
	UnitId           *uuid.UUID              `json:"unit_id"`
	IsAvailable      *bool                   `json:"is_available"`
	IsActive         *bool                   `json:"is_active"`
	Variants         []ProductVariantRequest `json:"variants"`
}

type ProductUpdateRequest struct {
	BusinessId       uuid.UUID                     `json:"business_id" validate:"required"`
	CategoryId       *uuid.UUID                    `json:"category_id"`
	Name             string                        `json:"name" validate:"required"`
	HasVariant       bool                          `json:"has_variant"`
	Brand            *string                       `json:"brand"`
	Description      *string                       `json:"description"`
	Image            *string                       `json:"image"`
	BasePrice        *float64                      `json:"base_price"`
	SellPrice        *float64                      `json:"sell_price"`
	SKU              *string                       `json:"sku"`
	Stock            *int                          `json:"stock"`
	TrackStock       *bool                         `json:"track_stock"`
	IgnoreStockCheck *bool                         `json:"ignore_stock_check"`
	MinimumSales     *int                          `json:"minimum_sales"`
	DiscountId       *uuid.UUID                    `json:"discount_id"`
	BrandId          *uuid.UUID                    `json:"brand_id"`
	TaxId            *uuid.UUID                    `json:"tax_id"`
	UnitId           *uuid.UUID                    `json:"unit_id"`
	IsAvailable      *bool                         `json:"is_available"`
	IsActive         *bool                         `json:"is_active"`
	Variants         []ProductVariantUpdateRequest `json:"variants"`
}

type CategoryRequest struct {
	Name       string     `json:"name" validate:"required"`
	BusinessId uuid.UUID  `json:"business_id" validate:"required"`
	ParentId   *uuid.UUID `json:"parent_id"`
}

type ProductVariantRequest struct {
	BusinessId       uuid.UUID `json:"business_id" validate:"required"`
	Name             string    `json:"name" validate:"required"`
	Description      *string   `json:"description"`
	BasePrice        *float64  `json:"base_price"`
	SellPrice        *float64  `json:"sell_price"`
	MinimumSales     *int      `json:"minimum_sales"`
	SKU              *string   `json:"sku"`
	Stock            *int      `json:"stock"`
	TrackStock       *bool     `json:"track_stock"`
	IgnoreStockCheck *bool     `json:"ignore_stock_check"`
	IsAvailable      *bool     `json:"is_available"`
	IsActive         *bool     `json:"is_active"`
}

type ProductVariantUpdateRequest struct {
	Id               uuid.UUID `json:"id" validate:"required"`
	BusinessId       uuid.UUID `json:"business_id" validate:"required"`
	Name             string    `json:"name" validate:"required"`
	Description      *string   `json:"description"`
	BasePrice        *float64  `json:"base_price" validate:"gte=0"`
	SellPrice        *float64  `json:"sell_price" validate:"gte=0"`
	MinimumSales     *int      `json:"minimum_sales"`
	SKU              *string   `json:"sku"`
	Stock            *int      `json:"stock"`
	IgnoreStockCheck *bool     `json:"ignore_stock_check"`
	TrackStock       *bool     `json:"track_stock"`
	IsAvailable      *bool     `json:"is_available"`
	IsActive         *bool     `json:"is_active"`
}
