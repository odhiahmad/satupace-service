package entity

import (
	"time"
)

type Business struct {
	Id             int    `gorm:"type:int;primary_key"`
	BusinessTypeId int    `gorm:"column:business_type_id;not null" json:"business_type_id"`
	Users          []User `gorm:"foreignKey:BusinessId"` // Many-to-Many
	Name           string `gorm:"type:varchar(255)" json:"name"`
	Email          string `gorm:"type:varchar(255) unique" json:"email"`
	OwnerName      string `gorm:"type:varchar(255)" json:"owner_name"`
	PhoneNumber    string `gorm:"uniqueIndex;type:varchar(255)" json:"phone_number"`
	Address        string `gorm:"type:varchar(255)" json:"address"`
	Lat            string `gorm:"type:varchar(255) null" json:"lat"`
	Long           string `gorm:"type:varchar(255) null" json:"long"`
	Logo           string `gorm:"type:varchar(255) null" json:"logo"`
	Rating         string `gorm:"type:varchar(255) null" json:"rating"`
	Image          string `gorm:"type:varchar(255) null" json:"image"`
	IsActive       bool   `gorm:"not null; column:is_active" json:"is_active"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
