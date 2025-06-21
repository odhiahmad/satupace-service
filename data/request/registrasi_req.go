package request

type RegistrationRequest struct {
	RoleId         int     `json:"role_id" validate:"required"`
	BusinessTypeId int     `json:"business_type_id" validate:"required"`
	Email          string  `json:"email" validate:"required,email"`
	Password       string  `json:"password" validate:"required,min=6"`
	Name           string  `json:"name" validate:"required"`
	OwnerName      string  `json:"owner_name" validate:"required"`
	Logo           *string `json:"logo,omitempty"`
	Rating         *string `json:"rating,omitempty"`
	Image          *string `json:"image,omitempty"`
}
