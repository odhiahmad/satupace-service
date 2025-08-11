package mapper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func MapTransaction(trx *entity.Transaction) *response.TransactionResponse {
	var itemResponses []response.TransactionItemResponse

	for _, item := range trx.Items {
		var attrResponses []response.TransactionItemAttributeResponse
		for _, attr := range item.Attributes {
			attrResponses = append(attrResponses, response.TransactionItemAttributeResponse{
				Id:                 attr.Id,
				ProductAttributeId: attr.ProductAttributeId,
				AdditionalPrice:    attr.AdditionalPrice,
			})
		}

		var productResponses []response.ProductResponse
		if item.Product != nil {
			productResponses = append(productResponses, MapProduct(*item.Product))
		}

		itemResponses = append(itemResponses, response.TransactionItemResponse{
			Id:                 item.Id,
			ProductId:          item.ProductId,
			Product:            productResponses,
			BundleId:           item.BundleId,
			ProductAttributeId: item.ProductAttributeId,
			ProductVariantId:   item.ProductVariantId,
			Quantity:           item.Quantity,
			BasePrice:          item.BasePrice,
			SellPrice:          item.SellPrice,
			Total:              item.Total,
			Discount:           item.Discount,
			Tax:                item.Tax,
			Promo:              item.Promo,
			Attributes:         attrResponses,
		})
	}

	var cashierRes response.UserBusinessResponse
	if trx.Cashier != nil {
		if cashierPtr := MapUserBusiness(*trx.Cashier); cashierPtr != nil {
			cashierRes = *cashierPtr
		}
	}

	return &response.TransactionResponse{
		Id:              trx.Id,
		BusinessId:      trx.BusinessId,
		CustomerId:      trx.CustomerId,
		Cashier:         cashierRes,
		PaymentMethodId: trx.PaymentMethodId,
		BillNumber:      trx.BillNumber,
		Items:           itemResponses,
		FinalPrice:      trx.FinalPrice,
		BasePrice:       trx.BasePrice,
		SellPrice:       trx.SellPrice,
		Discount:        trx.Discount,
		Promo:           trx.Promo,
		Tax:             trx.Tax,
		Status:          trx.Status,
		Rating:          trx.Rating,
		Notes:           trx.Notes,
		AmountReceived:  trx.AmountReceived,
		Change:          trx.Change,
		PaidAt:          trx.PaidAt,
		RefundedAt:      trx.RefundedAt,
		RefundedBy:      trx.RefundedBy,
		RefundReason:    trx.RefundReason,
		IsRefunded:      trx.IsRefunded,
		CanceledAt:      trx.CanceledAt,
		CanceledBy:      trx.CanceledBy,
		CanceledReason:  trx.CanceledReason,
		IsCanceled:      trx.IsCanceled,
		CreatedAt:       trx.CreatedAt,
		UpdatedAt:       trx.UpdatedAt,
	}
}
