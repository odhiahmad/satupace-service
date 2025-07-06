package entity

import "time"

type Category struct {
	Id         int        `gorm:"primaryKey" json:"id"`
	BusinessId int        `gorm:"not null" json:"business_id"`
	Business   *Business  `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Name       string     `gorm:"type:varchar(255)" json:"name"`
	ParentId   *int       `json:"parent_id,omitempty"`
	Children   []Category `gorm:"foreignKey:ParentId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"children,omitempty"`
	IsActive   bool       `gorm:"not null" json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
