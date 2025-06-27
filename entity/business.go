package entity

import "time"

type Business struct {
	Id             int              `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string           `gorm:"type:varchar(255);not null" json:"business_name"`
	OwnerName      string           `gorm:"type:varchar(255);not null" json:"owner_name"`
	BusinessTypeId int              `gorm:"not null" json:"business_type_id"`
	BusinessType   *BusinessType    `gorm:"foreignKey:BusinessTypeId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Branches       []BusinessBranch `gorm:"foreignKey:BusinessId" json:"branches,omitempty"`
	Image          *string          `gorm:"type:varchar(255)" json:"image,omitempty"`
	IsActive       bool             `gorm:"not null" json:"is_active"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}
