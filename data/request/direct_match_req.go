package request

type CreateDirectMatchRequest struct {
	User2Id string `json:"user_2_id" validate:"required"`
}

type UpdateDirectMatchStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending accepted rejected"`
}
