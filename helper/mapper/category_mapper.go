package mapper

import (
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
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
