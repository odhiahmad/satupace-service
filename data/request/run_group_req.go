package request

type CreateRunGroupRequest struct {
	Name              *string `json:"name"`
	AvgPace           float64 `json:"avg_pace" validate:"required"`
	PreferredDistance int     `json:"preferred_distance" validate:"required"`
	Latitude          float64 `json:"latitude" validate:"required"`
	Longitude         float64 `json:"longitude" validate:"required"`
	ScheduledAt       string  `json:"scheduled_at" validate:"required"`
	MaxMember         int     `json:"max_member" validate:"required"`
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
