package request

type CreateRunGroupScheduleRequest struct {
	DayOfWeek int    `json:"day_of_week" binding:"gte=0,lte=6"` // 0=Minggu s/d 6=Sabtu
	StartTime string `json:"start_time" binding:"required"`      // format "HH:MM"
}

type UpdateRunGroupScheduleRequest struct {
	DayOfWeek *int    `json:"day_of_week"`
	StartTime *string `json:"start_time"`
	IsActive  *bool   `json:"is_active"`
}
