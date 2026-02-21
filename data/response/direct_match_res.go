package response

import "time"

type DirectMatchResponse struct {
	Id        string     `json:"id"`
	User1Id   string     `json:"user_1_id"`
	User2Id   string     `json:"user_2_id"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	MatchedAt *time.Time `json:"matched_at,omitempty"`
}

type DirectMatchDetailResponse struct {
	Id        string        `json:"id"`
	User1Id   string        `json:"user_1_id"`
	User1     *UserResponse `json:"user_1,omitempty"`
	User2Id   string        `json:"user_2_id"`
	User2     *UserResponse `json:"user_2,omitempty"`
	Status    string        `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	MatchedAt *time.Time    `json:"matched_at,omitempty"`
}
