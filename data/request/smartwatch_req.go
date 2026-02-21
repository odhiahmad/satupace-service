package request

type ConnectDeviceRequest struct {
	DeviceType   string `json:"device_type" binding:"required"` // garmin, apple_watch, fitbit, samsung, strava, suunto
	DeviceName   string `json:"device_name"`
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token"`
	ExternalId   string `json:"external_id"`
}

type DisconnectDeviceRequest struct {
	DeviceId string `json:"device_id" binding:"required"`
}

type SyncActivityRequest struct {
	DeviceId       string  `json:"device_id" binding:"required"`
	ExternalId     string  `json:"external_id" binding:"required"` // Prevent duplicates
	Distance       float64 `json:"distance" binding:"required"`
	Duration       int     `json:"duration" binding:"required"`
	AvgPace        float64 `json:"avg_pace"`
	MaxPace        float64 `json:"max_pace"`
	AvgHeartRate   int     `json:"avg_heart_rate"`
	MaxHeartRate   int     `json:"max_heart_rate"`
	Calories       int     `json:"calories"`
	Cadence        int     `json:"cadence"`
	ElevationGain  float64 `json:"elevation_gain"`
	StartLatitude  float64 `json:"start_latitude"`
	StartLongitude float64 `json:"start_longitude"`
	EndLatitude    float64 `json:"end_latitude"`
	EndLongitude   float64 `json:"end_longitude"`
	RouteData      string  `json:"route_data"`                       // JSON array of GPS coordinates
	RawData        string  `json:"raw_data"`                         // Raw JSON from device
	ActivityDate   string  `json:"activity_date" binding:"required"` // RFC3339
}

type BatchSyncRequest struct {
	DeviceId   string                `json:"device_id" binding:"required"`
	Activities []SyncActivityRequest `json:"activities" binding:"required,min=1,max=50"`
}
