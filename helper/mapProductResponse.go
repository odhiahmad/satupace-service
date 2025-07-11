package helper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

// SafeString digunakan untuk menghindari panic jika pointer string nil
func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func MapProductToResponse(product entity.Product) response.ProductResponse {
	return response.ProductResponse{
		Id:           product.Id,
		Name:         product.Name,
		Description:  product.Description,
		Image:        product.Image,
		BasePrice:    product.BasePrice,
		SellPrice:    product.SellPrice,
		FinalPrice:   Float64Ptr(CalculateFinalPrice(&product)),
		SKU:          product.SKU,
		Stock:        product.Stock,
		TrackStock:   product.TrackStock,
		MinimumSales: product.MinimumSales,
		IsAvailable:  product.IsAvailable,
		IsActive:     product.IsActive,
		HasVariant:   product.HasVariant,
		Variants:     MapProductVariants(product.Variants),
		Brand:        MapBrand(product.Brand),
		Category:     MapCategory(product.Category),
		Tax:          MapTax(product.Tax),
		Discount:     MapDiscount(product.Discount),
		Unit:         MapUnit(product.Unit),
	}
}

func MapProductVariants(variants []entity.ProductVariant) []response.ProductVariantResponse {
	var result []response.ProductVariantResponse
	for _, v := range variants {
		result = append(result, response.ProductVariantResponse{
			Id:        v.Id,
			Name:      v.Name,
			BasePrice: v.BasePrice,
			SKU:       v.SKU,
		})
	}
	return result
}

func MapCategory(category *entity.Category) *response.CategoryResponse {
	if category == nil || category.Id == 0 {
		return nil
	}

	return &response.CategoryResponse{
		Id:   category.Id,
		Name: category.Name,
	}
}

func MapTax(tax *entity.Tax) *response.TaxResponse {
	if tax == nil {
		return nil
	}

	return &response.TaxResponse{
		Id:           tax.Id,
		Name:         tax.Name,
		Amount:       tax.Amount,
		IsPercentage: tax.IsPercentage,
	}
}

func MapDiscount(discount *entity.Discount) *response.DiscountResponse {
	if discount == nil {
		return nil
	}

	return &response.DiscountResponse{
		Id:           discount.Id,
		Name:         discount.Name,
		Amount:       discount.Amount,
		IsPercentage: discount.IsPercentage,
		IsMultiple:   discount.IsMultiple,
		IsGlobal:     discount.IsGlobal,
		StartAt:      discount.StartAt,
		EndAt:        discount.EndAt,
	}
}

func MapUnit(unit *entity.Unit) *response.UnitResponse {
	if unit == nil {
		return nil
	}

	return &response.UnitResponse{
		Id:         unit.Id,
		Name:       unit.Name,
		Alias:      unit.Alias,
		Multiplier: unit.Multiplier,
	}
}

func MapBrand(unit *entity.Brand) *response.BrandResponse {
	if unit == nil {
		return nil
	}

	return &response.BrandResponse{
		Id:   unit.Id,
		Name: unit.Name,
	}
}

func CalculateFinalPrice(p *entity.Product) float64 {
	if p.SellPrice == nil {
		return 0
	}

	price := *p.SellPrice

	if p.Discount != nil {
		if p.Discount.IsPercentage {
			price -= price * p.Discount.Amount / 100
		} else {
			price -= p.Discount.Amount
		}
	}

	if price < 0 {
		price = 0
	}

	return price
}
