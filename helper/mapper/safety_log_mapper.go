package mapper

import (
	"run-sync/data/response"
	"run-sync/entity"
)

func MapSafetyLog(l *entity.SafetyLog) *response.SafetyLogResponse {
	if l == nil {
		return nil
	}

	return &response.SafetyLogResponse{
		Id:        l.Id.String(),
		UserId:    l.UserId.String(),
		MatchId:   l.MatchId.String(),
		Status:    l.Status,
		Reason:    l.Reason,
		CreatedAt: l.CreatedAt,
	}
}

func MapSafetyLogDetail(l *entity.SafetyLog, user *entity.User) *response.SafetyLogDetailResponse {
	if l == nil {
		return nil
	}

	return &response.SafetyLogDetailResponse{
		Id:        l.Id.String(),
		UserId:    l.UserId.String(),
		User:      MapUser(user),
		MatchId:   l.MatchId.String(),
		Status:    l.Status,
		Reason:    l.Reason,
		CreatedAt: l.CreatedAt,
	}
}
