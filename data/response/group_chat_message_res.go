package response

import "time"

type GroupChatMessageResponse struct {
	Id        string    `json:"id"`
	GroupId   string    `json:"group_id"`
	SenderId  string    `json:"sender_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupChatMessageDetailResponse struct {
	Id         string        `json:"id"`
	GroupId    string        `json:"group_id"`
	SenderId   string        `json:"sender_id"`
	SenderName string        `json:"sender_name"`
	Sender     *UserResponse `json:"sender,omitempty"`
	Message    string        `json:"message"`
	CreatedAt  time.Time     `json:"created_at"`
}
