package request

type NotificationFilterRequest struct {
	IsRead *bool  `form:"is_read"` // filter: sudah dibaca atau belum
	Type   string `form:"type"`    // filter berdasarkan tipe notifikasi
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

type MarkNotificationReadRequest struct {
	Ids []string `json:"ids" binding:"required"` // list id notifikasi yang akan ditandai sudah dibaca
}

type RegisterDeviceTokenRequest struct {
	FCMToken string `json:"fcm_token" binding:"required"`
	Platform string `json:"platform" binding:"required,oneof=android ios web"`
}

type RemoveDeviceTokenRequest struct {
	FCMToken string `json:"fcm_token" binding:"required"`
}
