package entity

import "time"

type BusinessBranch struct {
	Id          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	BusinessId  int       `gorm:"not null" json:"business_id"` // Relasi ke bisnis utama
	Address     *string   `gorm:"type:text" json:"address,omitempty"`
	PhoneNumber *string   `gorm:"type:varchar(255)" json:"phone_number,omitempty"`
	Rating      *string   `gorm:"type:varchar(255)" json:"rating,omitempty"`
	Provinsi    *string   `gorm:"type:varchar(255)" json:"provinsi,omitempty"`
	Kota        *string   `gorm:"type:text" json:"kota,omitempty"`
	Kecamatan   *string   `gorm:"type:varchar(255)" json:"kecamatan,omitempty"`
	PostalCode  *string   `gorm:"type:varchar(255)" json:"postal_code,omitempty"`
	Phone       string    `gorm:"type:varchar(20)" json:"phone"` // No telp cabang (opsional)
	IsMain      bool      `gorm:"default:false" json:"is_main"`  // Apakah cabang utama
	IsActive    bool      `gorm:"default:true" json:"is_active"` // Status aktif
	Business    Business  `gorm:"foreignKey:BusinessId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
