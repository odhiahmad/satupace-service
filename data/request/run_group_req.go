package request

type CreateRunGroupRequest struct {
	Name              *string `json:"name"`
	AvgPace           float64 `json:"avg_pace" binding:"required"`
	PreferredDistance int     `json:"preferred_distance" binding:"required"`
	Latitude          float64 `json:"latitude" binding:"required"`
	Longitude         float64 `json:"longitude" binding:"required"`
	ScheduledAt       string  `json:"scheduled_at" binding:"required"`
	MaxMember         int     `json:"max_member" binding:"required"`
	IsWomenOnly       bool    `json:"is_women_only"`
}

type UpdateRunGroupRequest struct {
	Name              *string  `json:"name"`
	AvgPace           *float64 `json:"avg_pace"`
	PreferredDistance *int     `json:"preferred_distance"`
	Latitude          *float64 `json:"latitude"`
	Longitude         *float64 `json:"longitude"`
	ScheduledAt       *string  `json:"scheduled_at"`
	MaxMember         *int     `json:"max_member"`
	IsWomenOnly       *bool    `json:"is_women_only"`
	Status            *string  `json:"status"`
}
