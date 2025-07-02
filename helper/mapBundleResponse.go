package helper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func MapBundleToResponse(p entity.Bundle) response.BundleResponse {
	var items []response.BundleItemResponse
	for _, i := range p.Items {
		product := i.Product

		var basePrice float64
		if product.BasePrice != nil {
			basePrice = *product.BasePrice
		}

		var sku *string
		if product.SKU != nil {
			sku = product.SKU
		}

		items = append(items, response.BundleItemResponse{
			Id:          i.Id,
			ProductId:   i.ProductId,
			Name:        product.Name,
			Description: product.Description,
			Image:       product.Image,
			BasePrice:   basePrice,
			SKU:         sku,
			Quantity:    i.Quantity,
		})
	}

	description := StringOrDefault(p.Description, "")
	image := StringOrDefault(p.Image, "")

	return response.BundleResponse{
		Id:          p.Id,
		Name:        p.Name,
		Description: description,
		Image:       image,
		BasePrice:   p.BasePrice,
		Stock:       p.Stock,
		IsAvailable: p.IsAvailable,
		IsActive:    p.IsActive,
		Items:       items,
	}
}
