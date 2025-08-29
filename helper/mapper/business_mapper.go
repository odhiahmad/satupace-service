package mapper

import (
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
)

func MapBusiness(business *entity.Business) *response.BusinessResponse {

	if business == nil || business.Id == uuid.Nil {
		return nil
	}

	var membership *response.MembershipResponse

	if business.Membership != nil {
		membership = MapMembership(*business.Membership)
	}

	return &response.BusinessResponse{
		Id:           business.Id,
		Name:         business.Name,
		OwnerName:    business.OwnerName,
		BusinessType: MapBusinessType(business.BusinessType),
		Image:        business.Image,
		IsActive:     business.IsActive,
		Membership:   membership,
	}
}
