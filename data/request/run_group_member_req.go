package request

type CreateRunGroupMemberRequest struct {
	GroupId string `json:"group_id" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
	Status  string `json:"status" binding:"required"`
}

type UpdateRunGroupMemberRequest struct {
	Status *string `json:"status"`
}

type JoinRunGroupRequest struct {
	GroupId string `json:"group_id" binding:"required"`
}
