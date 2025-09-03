package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"
	"loka-kasir/helper"
)

func MapBundle(p entity.Bundle) response.BundleResponse {
	var items []response.BundleItemResponse
	for _, i := range p.Items {
		product := i.Product

		var basePrice float64
		var sellPrice float64

		if product.BasePrice != nil {
			basePrice = *product.BasePrice
		}

		if product.SellPrice != nil {
			sellPrice = *product.SellPrice
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
			BasePrice:   &basePrice,
			SellPrice:   &sellPrice,
			SKU:         sku,
			Quantity:    i.Quantity,
		})
	}

	description := helper.StringOrDefault(p.Description, "")
	image := helper.StringOrDefault(p.Image, "")

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
