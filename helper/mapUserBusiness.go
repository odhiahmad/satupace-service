package helper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func MapUserBusinessResponse(user entity.UserBusiness) *response.UserBusinessResponse {
	var membership *response.MembershipResponse
	var businessType *response.BusinessTypeResponse

	if user.Membership != nil {
		membership = MapMembershipResponse(*user.Membership)
	}

	if user.Business.BusinessType != nil {
		businessType = MapBusinessTypeToResponse(user.Business.BusinessType)
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
			Id:           user.Business.Id,
			Name:         user.Business.Name,
			OwnerName:    user.Business.OwnerName,
			BusinessType: businessType,
			Image:        user.Business.Image,
			IsActive:     user.Business.IsActive,
		},
		IsVerified: user.IsVerified,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
		Membership: membership,
	}
}

func MapMembershipResponse(m entity.Membership) *response.MembershipResponse {
	return &response.MembershipResponse{
		Id:        m.Id,
		Type:      m.Type,
		StartDate: m.StartDate,
		EndDate:   m.EndDate,
		IsActive:  m.IsActive,
	}
}

func MapBusinessTypeToResponse(m *entity.BusinessType) *response.BusinessTypeResponse {
	if m == nil {
		return nil
	}

	return &response.BusinessTypeResponse{
		Id:   m.Id,
		Name: m.Name,
	}
}
