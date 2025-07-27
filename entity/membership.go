package entity

import "time"

type Membership struct {
	Id        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    int       `gorm:"not null" json:"user_id"`
	StartDate time.Time `gorm:"not null" json:"start_date"`
	EndDate   time.Time `gorm:"not null" json:"end_date"`
	Type      string    `gorm:"type:varchar(20);not null" json:"type"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
