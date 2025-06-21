package helper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func ToProductResponse(p entity.Product) response.ProductResponse {
	var category *response.ProductCategoryResponse
	if p.ProductCategory.Id != 0 {
		category = &response.ProductCategoryResponse{
			Id:         p.ProductCategory.Id,
			Name:       p.ProductCategory.Name,
			BusinessId: p.ProductCategory.BusinessId,
			ParentId:   p.ProductCategory.ParentId,
		}
	}

	var variants []response.ProductVariantResponse
	for _, v := range p.Variants {
		variants = append(variants, response.ProductVariantResponse{
			Id:        v.Id,
			Name:      v.Name,
			BasePrice: v.BasePrice,
			SKU:       v.SKU,
			Stock:     v.Stock,
		})
	}

	return response.ProductResponse{
		Id:                p.Id,
		Name:              p.Name,
		Description:       p.Description,
		ProductCategoryId: p.ProductCategoryId,
		ProductCategory:   category,
		HasVariant:        p.HasVariant,
		BasePrice:         p.BasePrice,
		SKU:               p.SKU,
		Image:             p.Image,
		Variants:          variants,
	}
}
