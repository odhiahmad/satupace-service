package request

type SendDirectChatMessageRequest struct {
	MatchId string `json:"match_id" binding:"required"`
	Message string `json:"message" binding:"required"`
}
