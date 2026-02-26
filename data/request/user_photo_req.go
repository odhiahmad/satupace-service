package request

type UploadUserPhotoRequest struct {
	Image     string `json:"image" binding:"required"`
	Type      string `json:"type" binding:"required,oneof=profile run verification"`
	IsPrimary bool   `json:"is_primary"`
}

type UpdateUserPhotoRequest struct {
	Type      *string `json:"type"`
	IsPrimary *bool   `json:"is_primary"`
}

type FaceVerifyRequest struct {
	Image string `json:"image" binding:"required"` // base64 dari kamera
}
