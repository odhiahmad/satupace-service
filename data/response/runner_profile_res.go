package response

import "time"

type RunnerProfileResponse struct {
	Id                string    `json:"id"`
	UserId            string    `json:"user_id"`
	AvgPace           float64   `json:"avg_pace"`
	PreferredDistance int       `json:"preferred_distance"`
	PreferredTime     string    `json:"preferred_time"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	WomenOnlyMode     bool      `json:"women_only_mode"`
	Image             *string   `json:"image"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type RunnerProfileDetailResponse struct {
	Id                string        `json:"id"`
	UserId            string        `json:"user_id"`
	User              *UserResponse `json:"user,omitempty"`
	AvgPace           float64       `json:"avg_pace"`
	PreferredDistance int           `json:"preferred_distance"`
	PreferredTime     string        `json:"preferred_time"`
	Latitude          float64       `json:"latitude"`
	Longitude         float64       `json:"longitude"`
	WomenOnlyMode     bool          `json:"women_only_mode"`
	Image             *string       `json:"image"`
	IsActive          bool          `json:"is_active"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}
