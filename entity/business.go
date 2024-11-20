package entity

import (
	"time"

	"gorm.io/gorm"
)

type Business struct {
	gorm.Model
	BusinessBranch []BusinessBranch
	BusinessTypeID int
	BusinessType   BusinessType
	Name           string `gorm:"type:varchar(255)" json:"name"`
	OwnerName      string `gorm:"type:varchar(255)" json:"owner_name"`
	Logo           string `gorm:"type:varchar(255) null" json:"logo"`
	Rating         string `gorm:"type:varchar(255) null" json:"rating"`
	Image          string `gorm:"type:varchar(255) null" json:"image"`
	IsActive       bool   `gorm:"not null; column:is_active" json:"is_active"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
