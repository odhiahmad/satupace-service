package mapper

import (
	"run-sync/data/response"
	"run-sync/entity"
)

func MapGroupChatMessage(m *entity.GroupChatMessage) *response.GroupChatMessageResponse {
	if m == nil {
		return nil
	}

	return &response.GroupChatMessageResponse{
		Id:        m.Id.String(),
		GroupId:   m.GroupId.String(),
		SenderId:  m.SenderId.String(),
		Message:   m.Message,
		CreatedAt: m.CreatedAt,
	}
}

func MapGroupChatMessageDetail(m *entity.GroupChatMessage, sender *entity.User) *response.GroupChatMessageDetailResponse {
	if m == nil {
		return nil
	}

	return &response.GroupChatMessageDetailResponse{
		Id:        m.Id.String(),
		GroupId:   m.GroupId.String(),
		SenderId:  m.SenderId.String(),
		Sender:    MapUser(sender),
		Message:   m.Message,
		CreatedAt: m.CreatedAt,
	}
}
