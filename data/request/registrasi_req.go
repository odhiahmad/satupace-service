package request

type Registration struct {
	RoleId         int    `json:"role_id" validate:"required"`
	BusinessTypeId int    `json:"business_type_id" validate:"required"`
	Email          string `json:"email" validate:"required"`
	Password       string `json:"password" validate:"required"`
	Name           string `json:"name" validate:"required"`
	OwnerName      string `json:"owner_name" validate:"required"`
	Branch         []struct {
		Pic         string `json:"pic" validate:"required"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address" validate:"required"`
	}
}
