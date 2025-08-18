package request

type RegistrationRequest struct {
	Email          *string `json:"email"`
	Password       string  `json:"password" validate:"required,min=6"`
	BusinessTypeId int     `json:"business_type_id" validate:"required"`
	Name           string  `json:"name" validate:"required"`
	OwnerName      string  `json:"owner_name" validate:"required"`
	Address        *string `json:"address"`
	PhoneNumber    string  `json:"phone_number" validate:"required"`
	Type           *string `json:"type"` // hanya boleh 'monthly' atau 'yearly'
	Logo           *string `json:"logo"`
	Rating         *string `json:"rating"`
	Provinsi       *string `json:"provinsi"`
	Kota           *string `json:"kota"`
	Kecamatan      *string `json:"kecamatan"`
	PostalCode     *string `json:"postal_code"`
	Image          *string `json:"image"`
}
