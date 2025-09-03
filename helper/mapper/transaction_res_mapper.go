package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"
)

func MapTransactionItem(item entity.TransactionItem) response.TransactionItemResponse {
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

	var bundleResponses []response.BundleResponse
	if item.Bundle != nil {
		bundleResponses = append(bundleResponses, MapBundle(*item.Bundle))
	}

	return response.TransactionItemResponse{
		Id:                 item.Id,
		ProductId:          item.ProductId,
		Product:            productResponses,
		BundleId:           item.BundleId,
		Bundle:             bundleResponses,
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
	}
}

func MapTransaction(trx entity.Transaction) *response.TransactionResponse {
	var itemResponses []response.TransactionItemResponse
	for _, item := range trx.Items {
		itemResponses = append(itemResponses, MapTransactionItem(item))
	}

	var cashierRes response.UserBusinessResponse
	if trx.Cashier != nil {
		if cashierPtr := MapUserBusiness(*trx.Cashier); cashierPtr != nil {
			cashierRes = *cashierPtr
		}
	}

	var customerRes response.CustomerResponse
	if trx.Customer != nil {
		if customerPtr := MapCustomer(trx.Customer); customerPtr != nil {
			customerRes = *customerPtr
		}
	}

	return &response.TransactionResponse{
		Id:              trx.Id,
		BusinessId:      trx.BusinessId,
		Customer:        customerRes,
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
		OrderType:       *MapOrderType(trx.OrderType),
		Table:           MapTable(trx.Table),
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

func MapTransactions(trxs []entity.Transaction) []response.TransactionResponse {
	var result []response.TransactionResponse
	for _, trx := range trxs {
		result = append(result, *MapTransaction(trx))
	}
	return result
}

func MapTransactionItems(items []entity.TransactionItem) []response.TransactionItemResponse {
	var result []response.TransactionItemResponse
	for _, item := range items {
		result = append(result, MapTransactionItem(item))
	}
	return result
}
