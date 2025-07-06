package helper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func MapUserBusinessResponse(user entity.UserBusiness) *response.UserBusinessResponse {
	var memberships []response.MembershipResponse
	for _, m := range user.Memberships {
		memberships = append(memberships, response.MembershipResponse{
			Id:        m.Id,
			Type:      m.Type,
			StartDate: m.StartDate,
			EndDate:   m.EndDate,
			IsActive:  m.IsActive,
		})
	}

	return &response.UserBusinessResponse{
		Id:          user.Id,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role: &response.RoleResponse{
			Id:   user.Role.Id,
			Name: user.Role.Name,
		},
		Business: &response.BusinessResponse{
			Id:             user.Business.Id,
			Name:           user.Business.Name,
			OwnerName:      user.Business.OwnerName,
			BusinessTypeId: user.Business.BusinessTypeId,
			Image:          user.Business.Image,
			IsActive:       user.Business.IsActive,
		},
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Memberships: memberships,
	}
}
