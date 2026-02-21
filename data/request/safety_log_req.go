package request

type CreateSafetyLogRequest struct {
	MatchId string `json:"match_id" validate:"required"`
	Status  string `json:"status" validate:"required,oneof=reported blocked"`
	Reason  string `json:"reason" validate:"required"`
}
