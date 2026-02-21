package mapper

import (
	"run-sync/data/response"
	"run-sync/entity"
)

func MapRunGroup(g *entity.RunGroup) *response.RunGroupResponse {
	if g == nil {
		return nil
	}

	name := g.Name

	return &response.RunGroupResponse{
		Id:                g.Id.String(),
		Name:              name,
		AvgPace:           g.AvgPace,
		PreferredDistance: g.PreferredDistance,
		Latitude:          g.Latitude,
		Longitude:         g.Longitude,
		ScheduledAt:       g.ScheduledAt,
		MaxMember:         g.MaxMember,
		IsWomenOnly:       g.IsWomenOnly,
		Status:            g.Status,
		CreatedBy:         g.CreatedBy.String(),
		CreatedAt:         g.CreatedAt,
	}
}

func MapRunGroupDetail(g *entity.RunGroup, creator *entity.User, memberCount int) *response.RunGroupDetailResponse {
	if g == nil {
		return nil
	}

	name := g.Name

	return &response.RunGroupDetailResponse{
		Id:                g.Id.String(),
		Name:              name,
		AvgPace:           g.AvgPace,
		PreferredDistance: g.PreferredDistance,
		Latitude:          g.Latitude,
		Longitude:         g.Longitude,
		ScheduledAt:       g.ScheduledAt,
		MaxMember:         g.MaxMember,
		IsWomenOnly:       g.IsWomenOnly,
		Status:            g.Status,
		CreatedBy:         g.CreatedBy.String(),
		Creator:           MapUser(creator),
		MemberCount:       memberCount,
		CreatedAt:         g.CreatedAt,
	}
}
