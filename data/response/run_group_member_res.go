package response

import "time"

type RunGroupMemberResponse struct {
	Id       string    `json:"id"`
	GroupId  string    `json:"group_id"`
	UserId   string    `json:"user_id"`
	Status   string    `json:"status"`
	JoinedAt time.Time `json:"joined_at"`
}

type RunGroupMemberDetailResponse struct {
	Id       string        `json:"id"`
	GroupId  string        `json:"group_id"`
	UserId   string        `json:"user_id"`
	User     *UserResponse `json:"user,omitempty"`
	Status   string        `json:"status"`
	JoinedAt time.Time     `json:"joined_at"`
}
