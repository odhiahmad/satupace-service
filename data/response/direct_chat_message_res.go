package response

import "time"

type DirectChatMessageResponse struct {
	Id        string    `json:"id"`
	MatchId   string    `json:"match_id"`
	SenderId  string    `json:"sender_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type DirectChatMessageDetailResponse struct {
	Id         string        `json:"id"`
	MatchId    string        `json:"match_id"`
	SenderId   string        `json:"sender_id"`
	SenderName string        `json:"sender_name"`
	Sender     *UserResponse `json:"sender,omitempty"`
	Message    string        `json:"message"`
	CreatedAt  time.Time     `json:"created_at"`
}
