package mapper

import (
	"run-sync/data/response"
	"run-sync/entity"
)

func MapRunnerProfile(p *entity.RunnerProfile) *response.RunnerProfileResponse {
	if p == nil {
		return nil
	}

	return &response.RunnerProfileResponse{
		Id:                p.Id.String(),
		UserId:            p.UserId.String(),
		AvgPace:           p.AvgPace,
		PreferredDistance: p.PreferredDistance,
		PreferredTime:     p.PreferredTime,
		Latitude:          p.Latitude,
		Longitude:         p.Longitude,
		WomenOnlyMode:     p.WomenOnlyMode,
		Image:             p.Image,
		IsActive:          p.IsActive,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
}

func MapRunnerProfileDetail(p *entity.RunnerProfile) *response.RunnerProfileDetailResponse {
	if p == nil {
		return nil
	}

	return &response.RunnerProfileDetailResponse{
		Id:                p.Id.String(),
		UserId:            p.UserId.String(),
		User:              MapUser(p.User),
		AvgPace:           p.AvgPace,
		PreferredDistance: p.PreferredDistance,
		PreferredTime:     p.PreferredTime,
		Latitude:          p.Latitude,
		Longitude:         p.Longitude,
		WomenOnlyMode:     p.WomenOnlyMode,
		Image:             p.Image,
		IsActive:          p.IsActive,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}
}
