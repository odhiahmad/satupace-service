package dto

import uuid "github.com/satori/go.uuid"

type PerusahaanCreateDTO struct {
	Nama     string `json:"nama" form:"nama" binding:"required"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password,omitempty" form:"password,omitempty" binding:"required" validate:"min:6"`
}

type PerusahaanUpdateDTO struct {
	ID       uuid.UUID
	Nama     string `json:"nama" form:"nama" binding:"required"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password,omitempty" form:"password,omitempty" binding:"required" validate:"min:6"`
}
