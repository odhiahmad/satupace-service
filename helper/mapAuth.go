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
		Id:             user.Business.Id,
		Name:           user.Business.Name,
		OwnerName:      user.Business.OwnerName,
		BusinessTypeId: user.Business.BusinessTypeId,
		Image:          user.Business.Image,
		IsActive:       user.Business.IsActive,
		Type:           businessType,
	}

	var branch *response.BusinessBranchResponse
	if user.Branch != nil {
		branch = &response.BusinessBranchResponse{
			Id:          user.Branch.Id,
			PhoneNumber: &user.Branch.Phone,
		}
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
		Branch:      branch,
		Memberships: MapMembershipResponse(*user.Membership),
	}
}
