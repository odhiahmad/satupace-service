package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func MapOrderType(orderType *entity.OrderType) *response.OrderTypeResponse {
	caser := cases.Title(language.Indonesian)
	if orderType == nil {
		return nil
	}

	return &response.OrderTypeResponse{
		Id:   orderType.Id,
		Name: caser.String(orderType.Name),
	}
}
