package mapper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
)

func MapAuth(user *entity.UserBusiness, token string) *response.AuthResponse {

	return &response.AuthResponse{
		Id:          user.Id,
		Email:       helper.SafeString(user.Email),
		PhoneNumber: user.PhoneNumber,
		Token:       token,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Role:        *MapRole(&user.Role),
		Business:    *MapBusiness(&user.Business),
		Memberships: MapMembership(*user.Membership),
	}
}
