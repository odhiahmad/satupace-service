package mapper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func MapUserBusiness(user entity.UserBusiness) *response.UserBusinessResponse {
	var membership *response.MembershipResponse

	if user.Membership != nil {
		membership = MapMembership(*user.Membership)
	}

	return &response.UserBusinessResponse{
		Id:          user.Id,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role: &response.RoleResponse{
			Id:   user.Role.Id,
			Name: user.Role.Name,
		},
		Business:   MapBusiness(&user.Business),
		IsVerified: user.IsVerified,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		Membership: membership,
	}
}
