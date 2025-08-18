package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	Id         uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	BusinessId uuid.UUID      `gorm:"not null;index" json:"business_id"`
	Business   *Business      `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Name       string         `gorm:"type:varchar(100);not null" json:"name"`
	Phone      *string        `gorm:"type:varchar(20)" json:"phone"`
	Email      *string        `gorm:"type:varchar(100)" json:"email"`
	Address    *string        `gorm:"type:text" json:"address"`
	Notes      *string        `gorm:"type:text" json:"notes"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	if c.Id == uuid.Nil {
		c.Id = uuid.New()
	}
	return
}
