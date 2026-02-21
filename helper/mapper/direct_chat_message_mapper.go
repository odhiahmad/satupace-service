package mapper

import (
	"run-sync/data/response"
	"run-sync/entity"
)

func MapDirectChatMessage(m *entity.DirectChatMessage) *response.DirectChatMessageResponse {
	if m == nil {
		return nil
	}

	return &response.DirectChatMessageResponse{
		Id:        m.Id.String(),
		MatchId:   m.MatchId.String(),
		SenderId:  m.SenderId.String(),
		Message:   m.Message,
		CreatedAt: m.CreatedAt,
	}
}

func MapDirectChatMessageDetail(m *entity.DirectChatMessage, sender *entity.User) *response.DirectChatMessageDetailResponse {
	if m == nil {
		return nil
	}

	return &response.DirectChatMessageDetailResponse{
		Id:        m.Id.String(),
		MatchId:   m.MatchId.String(),
		SenderId:  m.SenderId.String(),
		Sender:    MapUser(sender),
		Message:   m.Message,
		CreatedAt: m.CreatedAt,
	}
}
