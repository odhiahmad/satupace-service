package helper

import (
	"time"

	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func MapProductToResponse(product entity.Product) response.ProductResponse {
	return response.ProductResponse{
		Id:               product.Id,
		Name:             product.Name,
		Description:      product.Description,
		Image:            product.Image,
		BasePrice:        product.BasePrice,
		SellPrice:        product.SellPrice,
		FinalPrice:       Float64Ptr(CalculateFinalPrice(&product)),
		SKU:              product.SKU,
		Stock:            product.Stock,
		TrackStock:       product.TrackStock,
		IgnoreStockCheck: product.IgnoreStockCheck,
		MinimumSales:     product.MinimumSales,
		IsAvailable:      product.IsAvailable != nil && *product.IsAvailable,
		IsActive:         product.IsActive != nil && *product.IsActive,
		HasVariant:       product.HasVariant,
		Variants:         MapProductVariants(product.Variants, &product),
		Brand:            MapBrand(product.Brand),
		Category:         MapCategory(product.Category),
		Tax:              MapTax(product.Tax),
		Discount:         MapDiscount(product.Discount),
		Unit:             MapUnit(product.Unit),
	}
}

func MapProductVariants(variants []entity.ProductVariant, product *entity.Product) []response.ProductVariantResponse {
	var result []response.ProductVariantResponse
	for _, v := range variants {
		result = append(result, response.ProductVariantResponse{
			Id:               v.Id,
			Name:             v.Name,
			Description:      v.Description,
			BasePrice:        v.BasePrice,
			SellPrice:        v.SellPrice,
			FinalPrice:       Float64Ptr(CalculateFinalPriceFromVariant(&v, product)),
			SKU:              v.SKU,
			Stock:            v.Stock,
			IgnoreStockCheck: v.IgnoreStockCheck,
			TrackStock:       v.TrackStock,
			IsAvailable:      v.IsAvailable != nil && *v.IsAvailable,
			IsActive:         v.IsActive != nil && *v.IsActive,
		})
	}
	return result
}

func MapCategory(category *entity.Category) *response.CategoryResponse {
	caser := cases.Title(language.Indonesian)

	if category == nil || category.Id == uuid.Nil {
		return nil
	}

	return &response.CategoryResponse{
		Id:   category.Id,
		Name: caser.String(category.Name),
	}
}

func MapTax(tax *entity.Tax) *response.TaxResponse {
	caser := cases.Title(language.Indonesian)
	if tax == nil {
		return nil
	}

	return &response.TaxResponse{
		Id:           tax.Id,
		Name:         caser.String(tax.Name),
		Amount:       tax.Amount,
		IsPercentage: tax.IsPercentage,
		IsGlobal:     tax.IsGlobal,
		IsActive:     tax.IsActive,
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
	caser := cases.Title(language.Indonesian)
	if unit == nil {
		return nil
	}

	return &response.UnitResponse{
		Id:    unit.Id,
		Name:  caser.String(unit.Name),
		Alias: caser.String(unit.Alias),
	}
}

func MapBrand(unit *entity.Brand) *response.BrandResponse {
	caser := cases.Title(language.Indonesian)
	if unit == nil {
		return nil
	}

	return &response.BrandResponse{
		Id:   unit.Id,
		Name: caser.String(unit.Name),
	}
}

func CalculateFinalPrice(p *entity.Product) float64 {
	if p.SellPrice == nil {
		return 0
	}

	price := *p.SellPrice

	if p.Discount != nil && IsDiscountActive(p.Discount) {
		if *p.Discount.IsPercentage {
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

func CalculateFinalPriceFromVariant(v *entity.ProductVariant, product *entity.Product) float64 {
	if v.SellPrice == nil {
		return 0
	}

	price := *v.SellPrice

	if product.Discount != nil && IsDiscountActive(product.Discount) {
		if *product.Discount.IsPercentage {
			price -= price * product.Discount.Amount / 100
		} else {
			price -= product.Discount.Amount
		}
	}

	if price < 0 {
		price = 0
	}

	return price
}

func IsDiscountActive(d *entity.Discount) bool {
	if !*d.IsActive {
		return false
	}

	now := time.Now()

	if !d.StartAt.IsZero() && now.Before(d.StartAt) {
		return false
	}

	if !d.EndAt.IsZero() && now.After(d.EndAt) {
		return false
	}

	return true
}
