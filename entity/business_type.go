package entity

import "time"

type BusinessType struct {
	Id        int       `gorm:"type:int;primary_key"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
