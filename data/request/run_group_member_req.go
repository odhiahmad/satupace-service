package request

type CreateRunGroupMemberRequest struct {
	GroupId string `json:"group_id" validate:"required"`
	UserId  string `json:"user_id" validate:"required"`
	Status  string `json:"status" validate:"required"`
}

type UpdateRunGroupMemberRequest struct {
	Status *string `json:"status"`
}

type JoinRunGroupRequest struct {
	GroupId string `json:"group_id" validate:"required"`
}
