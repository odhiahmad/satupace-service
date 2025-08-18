package mapper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
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
