package request

type UserUpdateDTO struct {
	ID       int
	Email    string `json:"email" form:"email"`
	Password string `json:"password,omitempty" form:"password,omitempty" binding:"required" validate:"min:6"`
}

type UserCreateDTO struct {
	Nama     string `json:"nama" form:"nama"`
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password,omitempty" form:"password,omitempty" binding:"required" validate:"min:6"`
}

type ProfileRequest struct {
	Id             int     `json:"id"`                                        // ID diperlukan untuk update                    // opsional                  // opsional                 // opsional
	Email          *string `json:"email,omitempty" binding:"omitempty,email"` // opsional dan valid email
	PhoneNumber    *string `json:"phone_number,omitempty"`                    // opsional                  // opsional                  // opsional
	BusinessName   *string `json:"business_name,omitempty"`
	OwnerName      *string `json:"owner_name,omitempty"`
	BusinessTypeId *int    `json:"business_type_id,omitempty"`
	Image          *string `json:"image,omitempty"`
	ProvinceId     *int    `json:"province_id,omitempty"`
	CityId         *int    `json:"city_id,omitempty"`
	DistrictId     *int    `json:"district_id,omitempty"`
	VillageId      *int    `json:"village_id,omitempty"`
}

type ChangePasswordRequest struct {
	Id          int    `json:"id"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ChangeEmailRequest struct {
	Id    int     `json:"id"`
	Email *string `json:"email" binding:"required,email"`
}

type ChangePhoneRequest struct {
	Id          int     `json:"id" binding:"required"`
	PhoneNumber *string `json:"phone_number" binding:"required"`
}
