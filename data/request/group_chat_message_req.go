package request

type SendGroupChatMessageRequest struct {
	GroupId string `json:"group_id" validate:"required"`
	Message string `json:"message" validate:"required"`
}
