package entity

import "time"

type BusinessType struct {
	Id          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"type:varchar(20);unique;not null" json:"code"`
	Name        string    `gorm:"type:varchar(255)" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
