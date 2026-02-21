package request

// StravaCallbackRequest is used when the frontend sends the OAuth authorization code.
type StravaCallbackRequest struct {
	Code string `json:"code" binding:"required"` // Authorization code from Strava OAuth redirect
}
