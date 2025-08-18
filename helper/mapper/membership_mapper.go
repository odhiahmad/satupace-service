package mapper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
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
