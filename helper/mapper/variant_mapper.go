package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"
	"loka-kasir/helper"
)

func MapVariants(variants []entity.ProductVariant, product *entity.Product) []response.ProductVariantResponse {
	var result []response.ProductVariantResponse
	for _, v := range variants {
		result = append(result, response.ProductVariantResponse{
			Id:               v.Id,
			Name:             v.Name,
			Description:      v.Description,
			BasePrice:        v.BasePrice,
			SellPrice:        v.SellPrice,
			FinalPrice:       helper.Float64Ptr(CalculateFinalPriceFromVariant(&v, product)),
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
