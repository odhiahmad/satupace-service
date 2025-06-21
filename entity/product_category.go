package entity

import "time"

type ProductCategory struct {
	Id          int      `gorm:"primaryKey"`
	BusinessId  int      `gorm:"not null"`
	Business    Business `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name        string   `gorm:"type:varchar(255)"`
	ParentId    *int
	Children    []ProductCategory `gorm:"foreignKey:ParentId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsAvailable bool              `gorm:"not null"`
	IsActive    bool              `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *ProductCategory) Prepare() error {
	p.IsActive = true
	return nil
}
