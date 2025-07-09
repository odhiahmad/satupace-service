package helper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func MapUserBusinessResponse(user entity.UserBusiness) *response.UserBusinessResponse {

	return &response.UserBusinessResponse{
		Id:          user.Id,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role: &response.RoleResponse{
			Id:   user.Role.Id,
			Name: user.Role.Name,
		},
		Business: &response.BusinessResponse{
			Id:        user.Business.Id,
			Name:      user.Business.Name,
			OwnerName: user.Business.OwnerName,
			Type: response.BusinessTypeResponse{
				Id:   user.Business.BusinessType.Id,
				Name: user.Business.BusinessType.Name,
			},
			Image:    user.Business.Image,
			IsActive: user.Business.IsActive,
		},
		IsVerified: user.IsVerified,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		Membership: MapMembershipResponse(*user.Membership),
	}
}
