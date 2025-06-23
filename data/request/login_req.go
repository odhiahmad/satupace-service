package request

type LoginUserBusinessDTO struct {
	Identifier string `json:"identifier" form:"identifier" binding:"required"`
	Password   string `json:"password,omitempty" form:"password,omitempty" validate:"min:5"`
}

type LoginUserDTO struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password,omitempty" form:"password,omitempty" validate:"min:5"`
}
