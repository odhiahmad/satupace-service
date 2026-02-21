package response

import "time"

type SafetyLogResponse struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	MatchId   string    `json:"match_id"`
	Status    string    `json:"status"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

type SafetyLogDetailResponse struct {
	Id        string        `json:"id"`
	UserId    string        `json:"user_id"`
	User      *UserResponse `json:"user,omitempty"`
	MatchId   string        `json:"match_id"`
	Status    string        `json:"status"`
	Reason    string        `json:"reason"`
	CreatedAt time.Time     `json:"created_at"`
}
