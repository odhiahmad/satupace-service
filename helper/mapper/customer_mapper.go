package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func MapCustomer(customer *entity.Customer) *response.CustomerResponse {
	caser := cases.Title(language.Indonesian)
	if customer == nil {
		return nil
	}

	return &response.CustomerResponse{
		Id:         customer.Id,
		BusinessId: customer.BusinessId,
		Name:       caser.String(customer.Name),
		Phone:      customer.Phone,
		Email:      customer.Email,
		Address:    customer.Address,
		CreatedAt:  customer.CreatedAt,
		UpdatedAt:  customer.UpdatedAt,
	}
}
