package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func MapBrand(unit *entity.Brand) *response.BrandResponse {
	caser := cases.Title(language.Indonesian)
	if unit == nil {
		return nil
	}

	return &response.BrandResponse{
		Id:   unit.Id,
		Name: caser.String(unit.Name),
	}
}
