package request

type UserUpdateDTO struct {
	ID       int
	Username string `json:"email" form:"email" binding:"required"`
	Password string `json:"password,omitempty" form:"password,omitempty" binding:"required" validate:"min:6"`
}

type UserCreateDTO struct {
	Nama     string `json:"nama" form:"nama" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password,omitempty" form:"password,omitempty" binding:"required" validate:"min:6"`
}
