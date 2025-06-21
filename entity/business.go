package entity

import (
	"time"
)

type Business struct {
	Id             int          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string       `gorm:"type:varchar(255);not null" json:"name"`
	OwnerName      string       `gorm:"type:varchar(255);not null" json:"owner_name"`
	BusinessTypeId int          `gorm:"not null" json:"business_type_id"`
	BusinessType   BusinessType `gorm:"foreignKey:BusinessTypeId" json:"-"`
	Logo           *string      `gorm:"type:varchar(255)" json:"logo,omitempty"`
	Rating         *string      `gorm:"type:varchar(255)" json:"rating,omitempty"`
	Image          *string      `gorm:"type:varchar(255)" json:"image,omitempty"`
	IsActive       bool         `gorm:"not null;column:is_active" json:"is_active"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

func (p *Business) Prepare() error {
	p.IsActive = true
	return nil
}
