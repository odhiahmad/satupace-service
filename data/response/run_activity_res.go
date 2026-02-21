package response

import "time"

type RunActivityResponse struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Distance  float64   `json:"distance"`
	Duration  int       `json:"duration"`
	AvgPace   float64   `json:"avg_pace"`
	Calories  int       `json:"calories"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}

type RunActivityDetailResponse struct {
	Id        string        `json:"id"`
	UserId    string        `json:"user_id"`
	User      *UserResponse `json:"user,omitempty"`
	Distance  float64       `json:"distance"`
	Duration  int           `json:"duration"`
	AvgPace   float64       `json:"avg_pace"`
	Calories  int           `json:"calories"`
	Source    string        `json:"source"`
	CreatedAt time.Time     `json:"created_at"`
}
