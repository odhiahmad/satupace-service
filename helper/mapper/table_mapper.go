package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func MapTable(table *entity.Table) *response.TableResponse {
	caser := cases.Title(language.Indonesian)
	if table == nil {
		return nil
	}

	return &response.TableResponse{
		Id:         table.Id,
		BusinessId: table.BusinessId,
		Number:     caser.String(table.Number),
		Status:     table.Status,
		CreatedAt:  table.CreatedAt,
		UpdatedAt:  table.UpdatedAt,
	}
}

func MapTableWithTransactions(tbl *entity.Table) *response.TableWithTransactionsResponse {
	var transactions []response.TransactionResponse
	for _, trx := range tbl.Transactions {
		transactions = append(transactions, response.TransactionResponse{
			Id:         trx.Id,
			Status:     trx.Status,
			Customer:   *MapCustomer(trx.Customer),
			Cashier:    *MapUserBusiness(*trx.Cashier),
			OrderType:  *MapOrderType(trx.OrderType),
			FinalPrice: trx.FinalPrice,
			CreatedAt:  trx.CreatedAt,
		})
	}

	return &response.TableWithTransactionsResponse{
		Id:           tbl.Id,
		Number:       tbl.Number,
		Status:       tbl.Status,
		Transactions: transactions,
	}
}
