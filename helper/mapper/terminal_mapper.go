package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"
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
