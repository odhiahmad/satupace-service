package mapper

import (
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/helper"
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

	senderName := ""
	if sender != nil {
		senderName = helper.DerefOrEmpty(sender.Name)
	}

	return &response.DirectChatMessageDetailResponse{
		Id:         m.Id.String(),
		MatchId:    m.MatchId.String(),
		SenderId:   m.SenderId.String(),
		SenderName: senderName,
		Sender:     MapUser(sender),
		Message:    m.Message,
		CreatedAt:  m.CreatedAt,
	}
}
