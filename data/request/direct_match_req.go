package request

type CreateDirectMatchRequest struct {
	User2Id string `json:"user_2_id" binding:"required"`
}

type UpdateDirectMatchStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending accepted rejected"`
}
