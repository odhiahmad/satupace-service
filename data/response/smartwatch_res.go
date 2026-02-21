package response

import "time"

type SmartWatchDeviceResponse struct {
	Id          string     `json:"id"`
	UserId      string     `json:"user_id"`
	DeviceType  string     `json:"device_type"`
	DeviceName  string     `json:"device_name"`
	ExternalId  string     `json:"external_id,omitempty"`
	IsConnected bool       `json:"is_connected"`
	LastSyncAt  *time.Time `json:"last_sync_at,omitempty"`
	ConnectedAt time.Time  `json:"connected_at"`
}

type SmartWatchSyncResponse struct {
	Id             string    `json:"id"`
	DeviceId       string    `json:"device_id"`
	UserId         string    `json:"user_id"`
	ActivityId     string    `json:"activity_id,omitempty"`
	ExternalId     string    `json:"external_id"`
	Distance       float64   `json:"distance"`
	Duration       int       `json:"duration"`
	AvgPace        float64   `json:"avg_pace"`
	MaxPace        float64   `json:"max_pace"`
	AvgHeartRate   int       `json:"avg_heart_rate"`
	MaxHeartRate   int       `json:"max_heart_rate"`
	Calories       int       `json:"calories"`
	Cadence        int       `json:"cadence"`
	ElevationGain  float64   `json:"elevation_gain"`
	StartLatitude  float64   `json:"start_latitude"`
	StartLongitude float64   `json:"start_longitude"`
	EndLatitude    float64   `json:"end_latitude"`
	EndLongitude   float64   `json:"end_longitude"`
	RouteData      string    `json:"route_data,omitempty"`
	Status         string    `json:"status"`
	ActivityDate   time.Time `json:"activity_date"`
	SyncedAt       time.Time `json:"synced_at"`
}

type SyncSummaryResponse struct {
	TotalSynced int                      `json:"total_synced"`
	TotalFailed int                      `json:"total_failed"`
	Activities  []SmartWatchSyncResponse `json:"activities"`
}

type DeviceStatsResponse struct {
	Device           SmartWatchDeviceResponse `json:"device"`
	TotalActivities  int                      `json:"total_activities"`
	TotalDistance    float64                  `json:"total_distance_km"`
	TotalDuration    int                      `json:"total_duration_seconds"`
	AvgPace          float64                  `json:"avg_pace"`
	AvgHeartRate     float64                  `json:"avg_heart_rate"`
	TotalCalories    int                      `json:"total_calories"`
	TotalElevation   float64                  `json:"total_elevation_m"`
	LastActivityDate *time.Time               `json:"last_activity_date,omitempty"`
}
