package mapper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func MapTax(tax *entity.Tax) *response.TaxResponse {
	caser := cases.Title(language.Indonesian)
	if tax == nil {
		return nil
	}

	return &response.TaxResponse{
		Id:           tax.Id,
		Name:         caser.String(tax.Name),
		Amount:       tax.Amount,
		IsPercentage: tax.IsPercentage,
		IsGlobal:     tax.IsGlobal,
		IsActive:     tax.IsActive,
	}
}
