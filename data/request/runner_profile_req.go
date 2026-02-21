package request

type CreateRunnerProfileRequest struct {
	AvgPace           float64 `json:"avg_pace" binding:"required"`
	PreferredDistance int     `json:"preferred_distance" binding:"required"`
	PreferredTime     string  `json:"preferred_time" binding:"required"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	WomenOnlyMode     bool    `json:"women_only_mode"`
	Image             *string `json:"image"`
}

type UpdateRunnerProfileRequest struct {
	AvgPace           *float64 `json:"avg_pace"`
	PreferredDistance *int     `json:"preferred_distance"`
	PreferredTime     *string  `json:"preferred_time"`
	Latitude          *float64 `json:"latitude"`
	Longitude         *float64 `json:"longitude"`
	WomenOnlyMode     *bool    `json:"women_only_mode"`
	Image             *string  `json:"image"`
}
