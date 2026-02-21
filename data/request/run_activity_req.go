package request

type CreateRunActivityRequest struct {
	Distance float64 `json:"distance" validate:"required"`
	Duration int     `json:"duration" validate:"required"`
	AvgPace  float64 `json:"avg_pace" validate:"required"`
	Calories int     `json:"calories"`
	Source   string  `json:"source" validate:"required"`
}

type UpdateRunActivityRequest struct {
	Distance *float64 `json:"distance"`
	Duration *int     `json:"duration"`
	AvgPace  *float64 `json:"avg_pace"`
	Calories *int     `json:"calories"`
	Source   *string  `json:"source"`
}
