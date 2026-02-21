package entity

import (
	"time"

	"github.com/google/uuid"
)

type DirectChatMessage struct {
	Id        uuid.UUID
	MatchId   uuid.UUID
	SenderId  uuid.UUID
	Message   string
	CreatedAt time.Time
}
