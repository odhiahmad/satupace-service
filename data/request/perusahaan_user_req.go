package request

type PerusahaanUserCreateDTO struct {
	Nama         string `json:"nama" form:"nama" binding:"required"`
	Username     string `json:"username" form:"username" binding:"required"`
	Password     string `json:"password,omitempty" form:"password,omitempty" binding:"required" validate:"min:6"`
	PerusahaanID int    `gorm:"null"`
}

type PerusahaanUserUpdateDTO struct {
	Id       int
	Nama     string `json:"nama" form:"nama" binding:"required"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password,omitempty" form:"password,omitempty" binding:"required" validate:"min:6"`
}
