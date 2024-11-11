package entity

import (
	"time"
)

type Perusahaan struct {
	Id        int    `gorm:"type:int;primary_key"`
	Nama      string `gorm:"type:varchar(255)"`
	Alamat    string `gorm:"type:varchar(255)"`
	Lat       string `gorm:"type:varchar(255)"`
	Long      string `gorm:"type:varchar(255)"`
	Logo      string `gorm:"type:varchar(255)"`
	Rating    string `gorm:"type:varchar(255)"`
	Gambar    string `gorm:"type:varchar(255)"`
	IsActive  bool   `gorm:"not null; column:is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
