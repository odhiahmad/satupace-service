package mapper

import (
	"run-sync/data/response"
	"run-sync/entity"
)

func MapUserPhoto(p *entity.UserPhoto) *response.UserPhotoResponse {
	if p == nil {
		return nil
	}

	return &response.UserPhotoResponse{
		Id:        p.Id.String(),
		UserId:    p.UserId.String(),
		Url:       p.Url,
		Type:      p.Type,
		IsPrimary: p.IsPrimary,
		CreatedAt: p.CreatedAt,
	}
}
