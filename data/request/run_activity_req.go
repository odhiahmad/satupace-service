package request

type CreateRunActivityRequest struct {
	Distance float64 `json:"distance" binding:"required"`
	Duration int     `json:"duration" binding:"required"`
	AvgPace  float64 `json:"avg_pace" binding:"required"`
	Calories int     `json:"calories"`
	Source   string  `json:"source" binding:"required"`
}

type UpdateRunActivityRequest struct {
	Distance *float64 `json:"distance"`
	Duration *int     `json:"duration"`
	AvgPace  *float64 `json:"avg_pace"`
	Calories *int     `json:"calories"`
	Source   *string  `json:"source"`
}
