package response

import "time"

type StravaConnectionResponse struct {
	Id          string     `json:"id"`
	AthleteId   int64      `json:"athlete_id"`
	Scope       string     `json:"scope"`
	IsConnected bool       `json:"is_connected"`
	LastSyncAt  *time.Time `json:"last_sync_at,omitempty"`
	ConnectedAt time.Time  `json:"connected_at"`
}

type StravaActivityResponse struct {
	Id               string    `json:"id"`
	StravaId         int64     `json:"strava_id"`
	RunActivityId    string    `json:"run_activity_id,omitempty"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	DistanceKm       float64   `json:"distance_km"`
	MovingTime       int       `json:"moving_time_seconds"`
	ElapsedTime      int       `json:"elapsed_time_seconds"`
	TotalElevation   float64   `json:"total_elevation_m"`
	AveragePace      float64   `json:"average_pace_min_km"`
	MaxSpeed         float64   `json:"max_speed_m_s"`
	AverageHeartrate float64   `json:"average_heartrate,omitempty"`
	MaxHeartrate     float64   `json:"max_heartrate,omitempty"`
	Calories         float64   `json:"calories"`
	StartDate        time.Time `json:"start_date"`
	MapPolyline      string    `json:"map_polyline,omitempty"`
	Status           string    `json:"status"`
	SyncedAt         time.Time `json:"synced_at"`
}

type StravaSyncSummaryResponse struct {
	TotalSynced  int                      `json:"total_synced"`
	TotalSkipped int                      `json:"total_skipped"`
	TotalFailed  int                      `json:"total_failed"`
	Activities   []StravaActivityResponse `json:"activities"`
}

type StravaStatsResponse struct {
	Connection      StravaConnectionResponse `json:"connection"`
	TotalActivities int                      `json:"total_activities"`
	TotalDistanceKm float64                  `json:"total_distance_km"`
	TotalDuration   int                      `json:"total_duration_seconds"`
	AvgPace         float64                  `json:"avg_pace_min_km"`
	AvgHeartrate    float64                  `json:"avg_heartrate"`
	TotalCalories   int                      `json:"total_calories"`
	TotalElevation  float64                  `json:"total_elevation_m"`
}

type StravaAuthURLResponse struct {
	AuthURL string `json:"auth_url"`
}
