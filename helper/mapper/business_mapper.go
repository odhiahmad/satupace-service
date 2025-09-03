package mapper

import (
	"loka-kasir/data/response"
	"loka-kasir/entity"

	"github.com/google/uuid"
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
