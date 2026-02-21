package response

import "time"

type ExploreRunnerResponse struct {
	UserId            string  `json:"user_id"`
	Name              *string `json:"name"`
	Gender            *string `json:"gender"`
	AvgPace           float64 `json:"avg_pace"`
	PreferredDistance int     `json:"preferred_distance"`
	PreferredTime     string  `json:"preferred_time"`
	Image             *string `json:"image,omitempty"`
	DistanceKm        float64 `json:"distance_km"` // Distance from requester
	WomenOnlyMode     bool    `json:"women_only_mode"`
}

type ExploreGroupResponse struct {
	GroupId           string    `json:"group_id"`
	Name              *string   `json:"name"`
	AvgPace           float64   `json:"avg_pace"`
	PreferredDistance int       `json:"preferred_distance"`
	ScheduledAt       time.Time `json:"scheduled_at"`
	MaxMember         int       `json:"max_member"`
	CurrentMembers    int       `json:"current_members"`
	IsWomenOnly       bool      `json:"is_women_only"`
	Status            string    `json:"status"`
	DistanceKm        float64   `json:"distance_km"`
	CreatedBy         string    `json:"created_by"`
}

type ExploreResponse struct {
	Runners []ExploreRunnerResponse `json:"runners,omitempty"`
	Groups  []ExploreGroupResponse  `json:"groups,omitempty"`
}
