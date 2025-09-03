package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"
)

func MapMembership(m entity.Membership) *response.MembershipResponse {
	return &response.MembershipResponse{
		Id:        m.Id,
		Type:      m.Type,
		StartDate: m.StartDate,
		EndDate:   m.EndDate,
		IsActive:  m.IsActive,
	}
}
