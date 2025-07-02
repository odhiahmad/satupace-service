package helper

import (
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
)

func ToTransactionItemRequests(items []entity.TransactionItem) []request.TransactionItemUpdate {
	var result []request.TransactionItemUpdate
	for _, item := range items {
		var attrReqs []request.TransactionItemAttributeUpdate
		for _, attr := range item.Attributes {
			attrReqs = append(attrReqs, request.TransactionItemAttributeUpdate{
				ProductAttributeId: attr.ProductAttributeId,
				AdditionalPrice:    attr.AdditionalPrice,
			})
		}

		result = append(result, request.TransactionItemUpdate{
			ProductId:          item.ProductId,
			BundleId:           item.BundleId,
			ProductVariantId:   item.ProductVariantId,
			ProductAttributeId: item.ProductAttributeId,
			Quantity:           item.Quantity,
			Total:              item.Total,
			Discount:           item.Discount,
			Promo:              item.Promo,
			Rating:             item.Rating,
			Attributes:         attrReqs,
		})
	}
	return result
}
