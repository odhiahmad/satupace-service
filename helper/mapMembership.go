package helper

import (
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func MapMembershipResponse(m entity.Membership) response.MembershipResponse {
	return response.MembershipResponse{
		Id:        m.Id,
		UserId:    m.UserId,
		Type:      m.Type,
		StartDate: m.StartDate,
		EndDate:   m.EndDate,
		IsActive:  m.IsActive,
	}
}
