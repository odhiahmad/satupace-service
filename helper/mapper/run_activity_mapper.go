package mapper

import (
	"run-sync/data/response"
	"run-sync/entity"
)

func MapRunActivity(a *entity.RunActivity) *response.RunActivityResponse {
	if a == nil {
		return nil
	}

	return &response.RunActivityResponse{
		Id:        a.Id.String(),
		UserId:    a.UserId.String(),
		Distance:  a.Distance,
		Duration:  a.Duration,
		AvgPace:   a.AvgPace,
		Calories:  a.Calories,
		Source:    a.Source,
		CreatedAt: a.CreatedAt,
	}
}

func MapRunActivityDetail(a *entity.RunActivity, user *entity.User) *response.RunActivityDetailResponse {
	if a == nil {
		return nil
	}

	return &response.RunActivityDetailResponse{
		Id:        a.Id.String(),
		UserId:    a.UserId.String(),
		User:      MapUser(user),
		Distance:  a.Distance,
		Duration:  a.Duration,
		AvgPace:   a.AvgPace,
		Calories:  a.Calories,
		Source:    a.Source,
		CreatedAt: a.CreatedAt,
	}
}
