package helper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func MapTransactionResponse(trx *entity.Transaction) *response.TransactionResponse {
	var itemResponses []response.TransactionItemResponse

	for _, item := range trx.Items {
		// Map atribut
		var attrResponses []response.TransactionItemAttributeResponse
		for _, attr := range item.Attributes {
			attrResponses = append(attrResponses, response.TransactionItemAttributeResponse{
				Id:                 attr.Id,
				ProductAttributeId: attr.ProductAttributeId,
				AdditionalPrice:    attr.AdditionalPrice,
			})
		}

		// Map product (jika ada)
		var productResponses []response.ProductResponse
		if item.Product != nil {
			productResponses = append(productResponses, MapProductToResponse(*item.Product))
		}

		// Map item
		itemResponses = append(itemResponses, response.TransactionItemResponse{
			Id:                 item.Id,
			ProductId:          item.ProductId,
			Product:            productResponses,
			BundleId:           item.BundleId,
			ProductAttributeId: item.ProductAttributeId,
			ProductVariantId:   item.ProductVariantId,
			Quantity:           item.Quantity,
			UnitPrice:          item.UnitPrice,
			Total:              item.Total,
			Discount:           item.Discount,
			Tax:                item.Tax,
			Promo:              item.Promo,
			Attributes:         attrResponses,
		})
	}

	return &response.TransactionResponse{
		Id:              trx.Id,
		BusinessId:      trx.BusinessId,
		CustomerId:      trx.CustomerId,
		PaymentMethodId: trx.PaymentMethodId,
		BillNumber:      trx.BillNumber,
		Items:           itemResponses,
		FinalPrice:      trx.FinalPrice,
		BasePrice:       trx.BasePrice,
		Discount:        trx.Discount,
		Promo:           trx.Promo,
		Tax:             trx.Tax,
		Status:          trx.Status,
		Rating:          trx.Rating,
		Notes:           trx.Notes,
		AmountReceived:  trx.AmountReceived,
		Change:          trx.Change,
	}
}
