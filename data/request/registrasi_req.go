package request

type Registration struct {
	RoleId         int    `json:"role_id" validate:"required"`
	BusinessTypeId int    `json:"business_type_id" validate:"required"`
	Email          string `json:"email" validate:"required"`
	Password       string `json:"password" validate:"min:6 required"`
	Name           string `json:"nama" validate:"required"`
	OwnerName      string `json:"owner_name" validate:"required"`
	PhoneNumber    string `json:"phone_number" validate:"required"`
	Address        string `json:"alamat" validate:"required"`
}
