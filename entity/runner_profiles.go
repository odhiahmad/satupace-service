package entity

import (
	"time"

	"github.com/google/uuid"
)

type RunnerProfile struct {
	Id uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	UserId uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	User   *User     `gorm:"constraint:OnDelete:CASCADE;"`

	AvgPace           float64 `gorm:"type:decimal(4,2)" json:"avg_pace"`      // contoh 5.45
	PreferredDistance int     `json:"preferred_distance"`                     // 5,10,21
	PreferredTime     string  `gorm:"type:varchar(50)" json:"preferred_time"` // morning/evening
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	WomenOnlyMode     bool    `json:"women_only_mode"`

	Image *string `gorm:"type:varchar(255)" json:"image"`

	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
