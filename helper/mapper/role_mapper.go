package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"
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
