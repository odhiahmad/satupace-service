package request

type CreateRunGroupRequest struct {
	Name              *string `json:"name"`
	MinPace           float64 `json:"min_pace" binding:"required"`
	MaxPace           float64 `json:"max_pace" binding:"required"`
	PreferredDistance int     `json:"preferred_distance" binding:"required"`
	Latitude          float64 `json:"latitude" binding:"required"`
	Longitude         float64 `json:"longitude" binding:"required"`
	ScheduledAt       string  `json:"scheduled_at" binding:"required"`
	MaxMember         int     `json:"max_member" binding:"required"`
	IsWomenOnly       bool    `json:"is_women_only"`
}

type UpdateRunGroupRequest struct {
	Name              *string  `json:"name"`
	MinPace           *float64 `json:"min_pace"`
	MaxPace           *float64 `json:"max_pace"`
	PreferredDistance *int     `json:"preferred_distance"`
	Latitude          *float64 `json:"latitude"`
	Longitude         *float64 `json:"longitude"`
	ScheduledAt       *string  `json:"scheduled_at"`
	MaxMember         *int     `json:"max_member"`
	IsWomenOnly       *bool    `json:"is_women_only"`
	Status            *string  `json:"status"`
}

type RunGroupFilterRequest struct {
	Status      string `form:"status"`
	WomenOnly   string `form:"women_only"`
	MinPace     string `form:"min_pace"`
	MaxPace     string `form:"max_pace"`
	MaxDistance string `form:"max_distance"`
	Latitude    string `form:"latitude"`
	Longitude   string `form:"longitude"`
	RadiusKm    string `form:"radius_km"`
}
