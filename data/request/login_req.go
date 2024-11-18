package request

type LoginDTO struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password,omitempty" form:"password,omitempty" validate:"min:5"`
}
