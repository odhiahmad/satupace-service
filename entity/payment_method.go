package entity

import "time"

type PaymentMethod struct {
	Id        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string    `gorm:"type:varchar(20);unique;not null" json:"code"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
