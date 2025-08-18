package mapper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func MapUnit(unit *entity.Unit) *response.UnitResponse {
	caser := cases.Title(language.Indonesian)
	if unit == nil {
		return nil
	}

	return &response.UnitResponse{
		Id:    unit.Id,
		Name:  caser.String(unit.Name),
		Alias: caser.String(unit.Alias),
	}
}
