package request

type UploadUserPhotoRequest struct {
	Url       string `json:"url" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=profile run verification"`
	IsPrimary bool   `json:"is_primary"`
}

type UpdateUserPhotoRequest struct {
	Type      *string `json:"type"`
	IsPrimary *bool   `json:"is_primary"`
}
