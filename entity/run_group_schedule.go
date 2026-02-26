package entity

import (
	"time"

	"github.com/google/uuid"
)

type RunGroupSchedule struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	GroupId   uuid.UUID `gorm:"type:uuid;not null;index"`
	DayOfWeek int       `gorm:"not null"` // 0=Minggu, 1=Senin, 2=Selasa, 3=Rabu, 4=Kamis, 5=Jumat, 6=Sabtu
	StartTime string    `gorm:"type:varchar(5);not null"` // format "HH:MM", contoh "06:30"
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
