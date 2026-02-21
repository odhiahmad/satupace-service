package mapper

import (
	"run-sync/data/response"
	"run-sync/entity"
	"time"
)

func MapDirectMatch(m *entity.DirectMatch) *response.DirectMatchResponse {
	if m == nil {
		return nil
	}

	var matchedAt *time.Time
	if m.MatchedAt != nil {
		matchedAt = m.MatchedAt
	}

	return &response.DirectMatchResponse{
		Id:        m.Id.String(),
		User1Id:   m.User1Id.String(),
		User2Id:   m.User2Id.String(),
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		MatchedAt: matchedAt,
	}
}

func MapDirectMatchDetail(m *entity.DirectMatch, u1 *entity.User, u2 *entity.User) *response.DirectMatchDetailResponse {
	if m == nil {
		return nil
	}

	return &response.DirectMatchDetailResponse{
		Id:        m.Id.String(),
		User1Id:   m.User1Id.String(),
		User1:     MapUser(u1),
		User2Id:   m.User2Id.String(),
		User2:     MapUser(u2),
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		MatchedAt: m.MatchedAt,
	}
}
