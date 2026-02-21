package mapper

import (
	"run-sync/data/response"
	"run-sync/entity"
)

func MapUser(u *entity.User) *response.UserResponse {
	if u == nil {
		return nil
	}

	return &response.UserResponse{
		Id:          u.Id.String(),
		Name:        u.Name,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Gender:      u.Gender,
		IsVerified:  u.IsVerified,
		IsActive:    u.IsActive,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func MapUserDetail(u *entity.User) *response.UserDetailResponse {
	if u == nil {
		return nil
	}
	return &response.UserDetailResponse{
		Id:          u.Id.String(),
		Name:        u.Name,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Gender:      u.Gender,
		IsVerified:  u.IsVerified,
		IsActive:    u.IsActive,
		Token:       u.Token,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
