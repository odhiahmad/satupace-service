package request

type SendDirectChatMessageRequest struct {
	MatchId string `json:"match_id" validate:"required"`
	Message string `json:"message" validate:"required"`
}
