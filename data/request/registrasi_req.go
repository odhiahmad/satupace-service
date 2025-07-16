package request

type RegistrationRequest struct {
	RoleId         int     `json:"role_id" validate:"required"`
	Email          *string `json:"email" validate:"email"`
	Password       string  `json:"password" validate:"required,min=6"`
	BusinessTypeId int     `json:"business_type_id" validate:"required"`
	Name           string  `json:"name" validate:"required"`
	OwnerName      string  `json:"owner_name" validate:"required"`
	Address        *string `json:"address,omitempty"`
	PhoneNumber    string  `json:"phone_number" validate:"required"`
	Type           *string `json:"type"` // hanya boleh 'monthly' atau 'yearly'
	Logo           *string `json:"logo,omitempty"`
	Rating         *string `json:"rating,omitempty"`
	Provinsi       *string `json:"provinsi,omitempty"`
	Kota           *string `json:"kota,omitempty"`
	Kecamatan      *string `json:"kecamatan,omitempty"`
	PostalCode     *string `json:"postal_code,omitempty"`
	Image          *string `json:"image,omitempty"`
}
