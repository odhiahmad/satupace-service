package helper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func ToDiscountResponse(discount entity.Discount) response.DiscountResponse {
	return response.DiscountResponse{
		Id:           discount.Id,
		Name:         discount.Name,
		Description:  discount.Description,
		IsPercentage: discount.IsPercentage,
		Amount:       discount.Amount,
		IsGlobal:     discount.IsGlobal,
		IsMultiple:   discount.IsMultiple,
		StartAt:      discount.StartAt,
		EndAt:        discount.EndAt,
		IsActive:     discount.IsActive,
	}
}
