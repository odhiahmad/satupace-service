package entity

import (
	"time"

	"github.com/google/uuid"
)

type RunGroup struct {
	Id uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	Name              *string `gorm:"type:varchar(255)"` // optional auto-generated
	MinPace           float64 `gorm:"type:decimal(4,2)"`
	MaxPace           float64 `gorm:"type:decimal(4,2)"`
	PreferredDistance int
	Latitude          float64
	Longitude         float64

	ScheduledAt time.Time
	MaxMember   int
	IsWomenOnly bool

	Status string `gorm:"type:varchar(50)"` // open, full, completed

	CreatedBy uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
