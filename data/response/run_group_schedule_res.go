package response

import "time"

type RunGroupScheduleResponse struct {
	Id        string    `json:"id"`
	GroupId   string    `json:"group_id"`
	DayOfWeek int       `json:"day_of_week"` // 0=Minggu, 1=Senin, ..., 6=Sabtu
	DayName   string    `json:"day_name"`    // nama hari dalam Bahasa Indonesia
	StartTime string    `json:"start_time"`  // "HH:MM"
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
