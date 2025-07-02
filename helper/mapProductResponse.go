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

func ToProductToResponse(product entity.Product) response.ProductResponse {
	var promos []response.PromoResponse

	for _, pp := range product.ProductPromos {
		if pp.Promo != nil && pp.Promo.Id != 0 {
			p := pp.Promo

			// Konversi required products
			var requiredProducts []response.RequiredProductData
			for _, rp := range p.RequiredProducts {
				requiredProducts = append(requiredProducts, response.RequiredProductData{
					Id:   rp.Id,
					Name: rp.Name,
				})
			}

			// Bangun response
			promos = append(promos, response.PromoResponse{
				Id:               p.Id,
				BusinessId:       p.BusinessId,
				Name:             p.Name,
				Description:      p.Description,
				Type:             p.Type,
				Amount:           p.Amount,
				IsPercentage:     p.IsPercentage,
				MinSpend:         p.MinSpend,
				MinQuantity:      p.MinQuantity, // helper
				FreeProduct:      nil,           // isi kalau ada relasi free product
				RequiredProducts: requiredProducts,
				StartDate:        p.StartDate,
				EndDate:          p.EndDate,
				IsActive:         p.IsActive,
			})
		}
	}

	var categoryRes *response.ProductCategoryResponse
	if product.ProductCategory != nil && product.ProductCategory.Id != 0 {
		categoryRes = &response.ProductCategoryResponse{
			Id:   product.ProductCategory.Id,
			Name: product.ProductCategory.Name,
		}
	}

	var taxRes *response.TaxResponse
	if product.Tax != nil {
		taxRes = &response.TaxResponse{
			Id:           product.Tax.Id,
			Name:         product.Tax.Name,
			Amount:       product.Tax.Amount,
			IsPercentage: product.Tax.IsPercentage,
		}
	}

	var discountRes *response.DiscountResponse
	if product.Discount != nil {
		discountRes = &response.DiscountResponse{
			Id:           product.Discount.Id,
			Name:         product.Discount.Name,
			Amount:       product.Discount.Amount,
			IsPercentage: product.Discount.IsPercentage,
			IsMultiple:   product.Discount.IsMultiple,
			IsGlobal:     product.Discount.IsGlobal,
			StartAt:      product.Discount.StartAt,
			EndAt:        product.Discount.EndAt,
		}
	}

	var unitRes *response.UnitResponse
	if product.Unit != nil {
		unitRes = &response.UnitResponse{
			Id:         product.Unit.Id,
			Name:       product.Unit.Name,
			Alias:      product.Unit.Alias,
			Multiplier: product.Unit.Multiplier,
		}
	}

	var variants []response.ProductVariantResponse
	for _, variant := range product.Variants {
		// Mapping relasi variant

		variants = append(variants, response.ProductVariantResponse{
			Id:        variant.Id,
			Name:      variant.Name,
			BasePrice: variant.BasePrice,
			SKU:       variant.SKU,
		})
	}

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
		Variants:        variants,
		ProductCategory: categoryRes,
		Tax:             taxRes,
		Discount:        discountRes,
		Unit:            unitRes,
		Promos:          promos,
	}
}
