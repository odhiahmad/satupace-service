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
		Id:              product.Id,
		Name:            product.Name,
		Description:     product.Description,
		Image:           product.Image,
		BasePrice:       product.BasePrice,
		SKU:             product.SKU,
		Stock:           product.Stock,
		MinimumSales:    product.MinimumSales,
		IsAvailable:     product.IsAvailable,
		IsActive:        product.IsActive,
		HasVariant:      product.HasVariant,
		Variants:        MapProductVariants(product.Variants),
		ProductCategory: MapProductCategory(product.ProductCategory),
		Tax:             MapTax(product.Tax),
		Discount:        MapDiscount(product.Discount),
		Unit:            MapUnit(product.Unit),
		Promos:          MapProductPromos(product.ProductPromos),
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

func MapProductCategory(category *entity.ProductCategory) *response.ProductCategoryResponse {
	if category == nil || category.Id == 0 {
		return nil
	}

	return &response.ProductCategoryResponse{
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
		BusinessId:   tax.BusinessId,
		Name:         tax.Name,
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

func MapProductPromos(productPromos []entity.ProductPromo) []response.PromoResponse {
	var result []response.PromoResponse

	for _, pp := range productPromos {
		if pp.Promo != nil && pp.Promo.Id != 0 {
			p := pp.Promo

			var requiredProducts []response.RequiredProductData
			for _, rp := range p.RequiredProducts {
				requiredProducts = append(requiredProducts, response.RequiredProductData{
					Id:   rp.Id,
					Name: rp.Name,
				})
			}

			result = append(result, response.PromoResponse{
				Id:               p.Id,
				BusinessId:       p.BusinessId,
				Name:             p.Name,
				Description:      p.Description,
				Type:             p.Type,
				Amount:           p.Amount,
				IsPercentage:     p.IsPercentage,
				MinSpend:         p.MinSpend,
				MinQuantity:      p.MinQuantity,
				FreeProduct:      nil, // tambahkan mapping jika perlu
				RequiredProducts: requiredProducts,
				StartDate:        p.StartDate,
				EndDate:          p.EndDate,
				IsActive:         p.IsActive,
			})
		}
	}

	return result
}
