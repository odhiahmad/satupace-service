package mapper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func MapTerminal(t *entity.Terminal) *response.TerminalResponse {
	return &response.TerminalResponse{
		Id:         t.Id,
		BusinessId: t.BusinessId,
		Name:       t.Name,
		Location:   t.Location,
		IsActive:   *t.IsActive,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}
