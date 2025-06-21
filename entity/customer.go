package entity

import (
	"time"
)

type Customer struct {
	Id        int    `gorm:"type:int;primary_key"`
	Name      string `gorm:"type:varchar(255)" json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
