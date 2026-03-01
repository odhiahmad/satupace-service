package response

import "time"

type UserPhotoResponse struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Url       string    `json:"url"`
	Type      string    `json:"type"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
	Warning   string    `json:"warning,omitempty"`
}

type FaceVerifyResponse struct {
	Matched    bool               `json:"matched"`
	Similarity float32            `json:"similarity"`
	Photo      *UserPhotoResponse `json:"photo,omitempty"` // diisi jika wajah cocok
}
