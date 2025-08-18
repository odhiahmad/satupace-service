package request

import "github.com/google/uuid"

type ProfileRequest struct {
	Id             uuid.UUID `json:"id"`
	Email          *string   `json:"email"`
	PhoneNumber    *string   `json:"phone_number"`
	BusinessName   *string   `json:"business_name"`
	OwnerName      *string   `json:"owner_name"`
	BusinessTypeId *int      `json:"business_type_id"`
	Image          *string   `json:"image"`
	ProvinceId     *int      `json:"province_id"`
	CityId         *int      `json:"city_id"`
	DistrictId     *int      `json:"district_id"`
	VillageId      *int      `json:"village_id"`
}

type ChangePasswordRequest struct {
	Id          uuid.UUID `json:"id"`
	OldPassword string    `json:"old_password" validate:"required"`
	NewPassword string    `json:"new_password" validate:"required,min=6"`
}

type ChangeEmailRequest struct {
	Id    uuid.UUID `json:"id"`
	Email *string   `json:"email" validate:"required,email"`
}

type ChangePhoneRequest struct {
	Id          uuid.UUID `json:"id" validate:"required"`
	PhoneNumber *string   `json:"phone_number" validate:"required"`
}
