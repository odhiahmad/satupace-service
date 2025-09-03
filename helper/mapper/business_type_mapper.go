package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"
)

func MapBusinessType(m *entity.BusinessType) *response.BusinessTypeResponse {
	if m == nil {
		return nil
	}

	return &response.BusinessTypeResponse{
		Id:   m.Id,
		Name: m.Name,
	}
}
