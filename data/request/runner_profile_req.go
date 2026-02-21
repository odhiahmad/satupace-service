package request

type CreateRunnerProfileRequest struct {
	AvgPace           float64 `json:"avg_pace" validate:"required"`
	PreferredDistance int     `json:"preferred_distance" validate:"required"`
	PreferredTime     string  `json:"preferred_time" validate:"required"`
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
