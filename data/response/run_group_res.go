package response

import "time"

type RunGroupResponse struct {
	Id                string    `json:"id"`
	Name              *string   `json:"name"`
	AvgPace           float64   `json:"avg_pace"`
	PreferredDistance int       `json:"preferred_distance"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	ScheduledAt       time.Time `json:"scheduled_at"`
	MaxMember         int       `json:"max_member"`
	IsWomenOnly       bool      `json:"is_women_only"`
	Status            string    `json:"status"`
	CreatedBy         string    `json:"created_by"`
	CreatedAt         time.Time `json:"created_at"`
}

type RunGroupDetailResponse struct {
	Id                string        `json:"id"`
	Name              *string       `json:"name"`
	AvgPace           float64       `json:"avg_pace"`
	PreferredDistance int           `json:"preferred_distance"`
	Latitude          float64       `json:"latitude"`
	Longitude         float64       `json:"longitude"`
	ScheduledAt       time.Time     `json:"scheduled_at"`
	MaxMember         int           `json:"max_member"`
	IsWomenOnly       bool          `json:"is_women_only"`
	Status            string        `json:"status"`
	CreatedBy         string        `json:"created_by"`
	Creator           *UserResponse `json:"creator,omitempty"`
	MemberCount       int           `json:"member_count"`
	CreatedAt         time.Time     `json:"created_at"`
}
