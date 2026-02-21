package request

type ExploreRunnersRequest struct {
	Latitude         float64 `form:"latitude" binding:"required"`
	Longitude        float64 `form:"longitude" binding:"required"`
	RadiusKm         float64 `form:"radius_km"`       // Default 10km
	MinPace          float64 `form:"min_pace"`        // min/km filter
	MaxPace          float64 `form:"max_pace"`        // min/km filter
	PreferredTime    string  `form:"preferred_time"`  // morning, evening, etc.
	Gender           string  `form:"gender"`          // male, female
	WomenOnly        bool    `form:"women_only"`      // Filter women-only mode
	Limit            int     `form:"limit"`           // Default 20
	ExcludeMatchedId bool    `form:"exclude_matched"` // Exclude already matched users
}

type ExploreGroupsRequest struct {
	Latitude  float64 `form:"latitude" binding:"required"`
	Longitude float64 `form:"longitude" binding:"required"`
	RadiusKm  float64 `form:"radius_km"` // Default 10km
	MinPace   float64 `form:"min_pace"`
	MaxPace   float64 `form:"max_pace"`
	WomenOnly bool    `form:"women_only"`
	Status    string  `form:"status"` // open (default), full, etc.
	Limit     int     `form:"limit"`  // Default 20
}
