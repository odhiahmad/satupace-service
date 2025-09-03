package mapper

import (
	"time"

	"loka-kasir/data/response"
	"loka-kasir/entity"
	"loka-kasir/helper"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func MapProduct(product entity.Product) response.ProductResponse {
	caser := cases.Title(language.Indonesian)

	return response.ProductResponse{
		Id:               product.Id,
		Name:             caser.String(product.Name),
		Description:      product.Description,
		Image:            product.Image,
		BasePrice:        product.BasePrice,
		SellPrice:        product.SellPrice,
		FinalPrice:       helper.Float64Ptr(CalculateFinalPrice(&product)),
		SKU:              product.SKU,
		Stock:            product.Stock,
		TrackStock:       product.TrackStock,
		IgnoreStockCheck: product.IgnoreStockCheck,
		MinimumSales:     product.MinimumSales,
		IsAvailable:      product.IsAvailable != nil && *product.IsAvailable,
		IsActive:         product.IsActive != nil && *product.IsActive,
		HasVariant:       product.HasVariant,
		Variants:         MapVariants(product.Variants, &product),
		Brand:            MapBrand(product.Brand),
		Category:         MapCategory(product.Category),
		Tax:              MapTax(product.Tax),
		Discount:         MapDiscount(product.Discount),
		Unit:             MapUnit(product.Unit),
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
