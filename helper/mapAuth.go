package helper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func MapAuthResponse(user *entity.UserBusiness, token string) *response.AuthResponse {
	role := response.RoleResponse{
		Id:   user.Role.Id,
		Name: user.Role.Name,
	}

	businessType := response.BusinessTypeResponse{
		Id:   user.Business.BusinessType.Id,
		Name: user.Business.BusinessType.Name,
	}

	business := response.BusinessResponse{
		Id:           user.Business.Id,
		Name:         user.Business.Name,
		OwnerName:    user.Business.OwnerName,
		Image:        user.Business.Image,
		IsActive:     user.Business.IsActive,
		BusinessType: &businessType,
	}

	return &response.AuthResponse{
		Id:          user.Id,
		Email:       *user.Email,
		PhoneNumber: user.PhoneNumber,
		Token:       token,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Role:        role,
		Business:    business,
		Memberships: MapMembershipResponse(*user.Membership),
	}
}
