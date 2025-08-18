package request

import (
	"github.com/google/uuid"
)

type TerminalRequest struct {
	BusinessId uuid.UUID `json:"business_id" binding:"required"`
	Name       string    `json:"name" binding:"required"`
	Location   string    `json:"location"`
	IsActive   *bool     `json:"is_active"`
}
