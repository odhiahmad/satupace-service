package mapper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
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
