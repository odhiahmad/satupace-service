package entity

import (
	"time"

	"github.com/google/uuid"
)

type SafetyLog struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	MatchId   uuid.UUID
	Status    string // reported, blocked
	Reason    string
	CreatedAt time.Time
}
