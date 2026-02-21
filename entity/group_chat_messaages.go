package entity

import (
	"time"

	"github.com/google/uuid"
)

type GroupChatMessage struct {
	Id        uuid.UUID
	GroupId   uuid.UUID
	SenderId  uuid.UUID
	Message   string
	CreatedAt time.Time
}
