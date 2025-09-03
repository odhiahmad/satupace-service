package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func MapDiscount(discount *entity.Discount) *response.DiscountResponse {
	caser := cases.Title(language.Indonesian)

	if discount == nil {
		return nil
	}

	return &response.DiscountResponse{
		Id:           discount.Id,
		Name:         caser.String(discount.Name),
		Description:  caser.String(discount.Description),
		IsPercentage: discount.IsPercentage,
		Amount:       discount.Amount,
		IsGlobal:     discount.IsGlobal,
		IsMultiple:   discount.IsMultiple,
		StartAt:      discount.StartAt,
		EndAt:        discount.EndAt,
		IsActive:     discount.IsActive,
	}
}
