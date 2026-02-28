package response

import "time"

type RunGroupResponse struct {
	Id                string                     `json:"id"`
	Name              *string                    `json:"name"`
	MinPace           float64                    `json:"min_pace"`
	MaxPace           float64                    `json:"max_pace"`
	PreferredDistance int                        `json:"preferred_distance"`
	DistanceKm        *float64                   `json:"distance_km,omitempty"`
	Latitude          float64                    `json:"latitude"`
	Longitude         float64                    `json:"longitude"`
	MeetingPoint      string                     `json:"meeting_point"`
	ScheduledAt       time.Time                  `json:"scheduled_at"`
	MaxMember         int                        `json:"max_member"`
	MemberCount       int                        `json:"member_count"`
	IsWomenOnly       bool                       `json:"is_women_only"`
	Status            string                     `json:"status"`
	CreatedBy         string                     `json:"created_by"`
	MyRole            string                     `json:"my_role,omitempty"`
	Schedules         []RunGroupScheduleResponse `json:"schedules"`
	CreatedAt         time.Time                  `json:"created_at"`
}

type RunGroupDetailResponse struct {
	Id                string                     `json:"id"`
	Name              *string                    `json:"name"`
	MinPace           float64                    `json:"min_pace"`
	MaxPace           float64                    `json:"max_pace"`
	PreferredDistance int                        `json:"preferred_distance"`
	Latitude          float64                    `json:"latitude"`
	Longitude         float64                    `json:"longitude"`
	MeetingPoint      string                     `json:"meeting_point"`
	ScheduledAt       time.Time                  `json:"scheduled_at"`
	MaxMember         int                        `json:"max_member"`
	IsWomenOnly       bool                       `json:"is_women_only"`
	Status            string                     `json:"status"`
	CreatedBy         string                     `json:"created_by"`
	Creator           *UserResponse              `json:"creator,omitempty"`
	MemberCount       int                        `json:"member_count"`
	Schedules         []RunGroupScheduleResponse `json:"schedules"`
	CreatedAt         time.Time                  `json:"created_at"`
}
