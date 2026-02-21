package entity

import (
	"time"

	"github.com/google/uuid"
)

// SmartWatchDevice represents a connected smartwatch/fitness tracker.
type SmartWatchDevice struct {
	Id           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId       uuid.UUID `gorm:"type:uuid;not null;index"`
	DeviceType   string    `gorm:"type:varchar(50);not null"` // garmin, apple_watch, fitbit, samsung, strava, suunto
	DeviceName   string    `gorm:"type:varchar(100)"`         // "Garmin Forerunner 265"
	AccessToken  string    `gorm:"type:text" json:"-"`
	RefreshToken string    `gorm:"type:text" json:"-"`
	ExternalId   string    `gorm:"type:varchar(255)"` // User/device ID on the external platform
	IsConnected  bool      `gorm:"default:true"`
	LastSyncAt   *time.Time
	ConnectedAt  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// SmartWatchSync represents a synced activity from a smartwatch.
type SmartWatchSync struct {
	Id             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DeviceId       uuid.UUID `gorm:"type:uuid;not null;index"`
	UserId         uuid.UUID `gorm:"type:uuid;not null;index"`
	ActivityId     uuid.UUID `gorm:"type:uuid;index"`               // Linked RunActivity
	ExternalId     string    `gorm:"type:varchar(255);uniqueIndex"` // External activity ID to prevent duplicates
	RawData        string    `gorm:"type:jsonb"`                    // Raw JSON from device API
	Distance       float64   // km
	Duration       int       // seconds
	AvgPace        float64   // min/km
	MaxPace        float64   // min/km
	AvgHeartRate   int
	MaxHeartRate   int
	Calories       int
	Cadence        int     // steps per minute
	ElevationGain  float64 // meters
	StartLatitude  float64
	StartLongitude float64
	EndLatitude    float64
	EndLongitude   float64
	RouteData      string `gorm:"type:jsonb"`                        // GPS coordinates array
	Status         string `gorm:"type:varchar(50);default:'synced'"` // synced, imported, failed
	SyncedAt       time.Time
	ActivityDate   time.Time // When the actual run happened
	CreatedAt      time.Time
}
