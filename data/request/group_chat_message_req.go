package request

type SendGroupChatMessageRequest struct {
	GroupId string `json:"group_id" binding:"required"`
	Message string `json:"message" binding:"required"`
}
