package entity

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Perusahaan struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Nama      string    `gorm:"type:varchar(255)" json:"nama"`
	Alamat    string    `gorm:"type:varchar(255)" json:"alamat"`
	Lat       string    `gorm:"type:varchar(255)" json:"lat"`
	Long      string    `gorm:"type:varchar(255)" json:"long"`
	Logo      string    `gorm:"type:varchar(255)" json:"logo"`
	Rating    string    `gorm:"type:varchar(255)" json:"rating"`
	Gambar    string    `gorm:"type:varchar(255)" json:"gambar"`
	IsActive  bool      `gorm:"not null; column:is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *Perusahaan) Prepare() error {
	u.ID = uuid.NewV4()
	u.IsActive = true
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}
