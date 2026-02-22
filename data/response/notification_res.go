package response

import "time"

type NotificationResponse struct {
	Id        string     `json:"id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	IsRead    bool       `json:"is_read"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	ActorId   *string    `json:"actor_id,omitempty"`
	RefId     *string    `json:"ref_id,omitempty"`
	RefType   *string    `json:"ref_type,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	UnreadCount   int64                  `json:"unread_count"`
}
