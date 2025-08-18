package mapper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func MapRole(role *entity.Role) *response.RoleResponse {
	if role == nil {
		return nil
	}

	return &response.RoleResponse{
		Id:   role.Id,
		Name: role.Name,
	}
}
