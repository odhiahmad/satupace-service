package entity

import (
	"time"

	"github.com/google/uuid"
)

// StravaConnection stores the OAuth connection between a user and their Strava account.
type StravaConnection struct {
	Id           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	AthleteId    int64     `gorm:"not null;uniqueIndex"` // Strava athlete ID
	AccessToken  string    `gorm:"type:text" json:"-"`   // Short-lived access token
	RefreshToken string    `gorm:"type:text" json:"-"`   // Long-lived refresh token
	ExpiresAt    int64     `gorm:"not null"`             // Unix timestamp when access token expires
	Scope        string    `gorm:"type:varchar(255)"`    // OAuth scope granted
	IsConnected  bool      `gorm:"default:true"`
	LastSyncAt   *time.Time
	ConnectedAt  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// StravaActivity stores synced activity data from Strava.
type StravaActivity struct {
	Id               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ConnectionId     uuid.UUID `gorm:"type:uuid;not null;index"`
	UserId           uuid.UUID `gorm:"type:uuid;not null;index"`
	RunActivityId    uuid.UUID `gorm:"type:uuid;index"`      // Linked RunActivity
	StravaId         int64     `gorm:"not null;uniqueIndex"` // Strava activity ID (prevent duplicates)
	Name             string    `gorm:"type:varchar(255)"`
	Type             string    `gorm:"type:varchar(50)"` // Run, Trail Run, etc.
	Distance         float64   // meters
	MovingTime       int       // seconds
	ElapsedTime      int       // seconds
	TotalElevation   float64   // meters
	AverageSpeed     float64   // meters/second
	MaxSpeed         float64   // meters/second
	AverageHeartrate float64
	MaxHeartrate     float64
	Calories         float64
	StartDate        time.Time
	StartLatitude    float64
	StartLongitude   float64
	EndLatitude      float64
	EndLongitude     float64
	MapPolyline      string `gorm:"type:text"`                         // Encoded polyline
	Status           string `gorm:"type:varchar(50);default:'synced'"` // synced, failed
	SyncedAt         time.Time
	CreatedAt        time.Time
}
