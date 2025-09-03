package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func MapCategory(category *entity.Category) *response.CategoryResponse {
	caser := cases.Title(language.Indonesian)

	if category == nil || category.Id == uuid.Nil {
		return nil
	}

	return &response.CategoryResponse{
		Id:   category.Id,
		Name: caser.String(category.Name),
	}
}
