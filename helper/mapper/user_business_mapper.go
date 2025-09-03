package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"
)

func MapUserBusiness(user entity.UserBusiness) *response.UserBusinessResponse {

	return &response.UserBusinessResponse{
		Id:          user.Id,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        MapRole(&user.Role),
		Business:    MapBusiness(&user.Business),
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
