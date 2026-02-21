package request

type CreateSafetyLogRequest struct {
	MatchId string `json:"match_id" binding:"required"`
	Status  string `json:"status" binding:"required,oneof=reported blocked"`
	Reason  string `json:"reason" binding:"required"`
}
