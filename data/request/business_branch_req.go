package request

type BusinessBranchCreate struct {
	BusinessId  int     `json:"business_id" validate:"required"`
	PhoneNumber *string `gorm:"type:varchar(255)" json:"phone_number,omitempty"`
	Rating      *string `gorm:"type:varchar(255)" json:"rating,omitempty"`
	Provinsi    *string `gorm:"type:varchar(255)" json:"provinsi,omitempty"`
	Kota        *string `gorm:"type:text" json:"kota,omitempty"`
	Kecamatan   *string `gorm:"type:varchar(255)" json:"kecamatan,omitempty"`
	PostalCode  *string `gorm:"type:varchar(255)" json:"postal_code,omitempty"`
	Phone       string  `gorm:"type:varchar(20)" json:"phone"` // No telp cabang (opsional)
	IsMain      bool    `json:"is_main"`
	IsActive    bool    `json:"is_active"`
}

type BusinessBranchUpdate struct {
	Id          int     `json:"id" validate:"required"`
	PhoneNumber *string `gorm:"type:varchar(255)" json:"phone_number,omitempty"`
	Rating      *string `gorm:"type:varchar(255)" json:"rating,omitempty"`
	Provinsi    *string `gorm:"type:varchar(255)" json:"provinsi,omitempty"`
	Kota        *string `gorm:"type:text" json:"kota,omitempty"`
	Kecamatan   *string `gorm:"type:varchar(255)" json:"kecamatan,omitempty"`
	PostalCode  *string `gorm:"type:varchar(255)" json:"postal_code,omitempty"`
	Phone       string  `gorm:"type:varchar(20)" json:"phone"` // No telp cabang (opsional)
	IsMain      bool    `json:"is_main"`
	IsActive    bool    `json:"is_active"`
}
