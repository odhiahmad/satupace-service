package entity

import (
	"time"

	"github.com/google/uuid"
)

type RunActivity struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId    uuid.UUID `gorm:"not null;index"`
	Distance  float64   // km
	Duration  int       // seconds
	AvgPace   float64
	Calories  int
	Source    string `gorm:"type:varchar(50)"` // manual, garmin, strava
	CreatedAt time.Time
}
