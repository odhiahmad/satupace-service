package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/repository"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StravaService interface {
	GetAuthURL(userId uuid.UUID) (response.StravaAuthURLResponse, error)
	HandleCallback(userId uuid.UUID, code string) (response.StravaConnectionResponse, error)
	Disconnect(userId uuid.UUID) error
	GetConnection(userId uuid.UUID) (response.StravaConnectionResponse, error)
	SyncActivities(userId uuid.UUID) (response.StravaSyncSummaryResponse, error)
	GetSyncHistory(userId uuid.UUID, limit int) ([]response.StravaActivityResponse, error)
	GetStats(userId uuid.UUID) (response.StravaStatsResponse, error)
}

type stravaService struct {
	stravaRepo   repository.StravaRepository
	activityRepo repository.RunActivityRepository
}

func NewStravaService(stravaRepo repository.StravaRepository, activityRepo repository.RunActivityRepository) StravaService {
	return &stravaService{stravaRepo: stravaRepo, activityRepo: activityRepo}
}

// ── Strava API structs ──

type stravaTokenResponse struct {
	TokenType    string        `json:"token_type"`
	ExpiresAt    int64         `json:"expires_at"`
	ExpiresIn    int           `json:"expires_in"`
	RefreshToken string        `json:"refresh_token"`
	AccessToken  string        `json:"access_token"`
	Athlete      stravaAthlete `json:"athlete"`
}

type stravaAthlete struct {
	Id        int64  `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type stravaActivityAPI struct {
	Id               int64             `json:"id"`
	Name             string            `json:"name"`
	Type             string            `json:"type"`
	SportType        string            `json:"sport_type"`
	Distance         float64           `json:"distance"`             // meters
	MovingTime       int               `json:"moving_time"`          // seconds
	ElapsedTime      int               `json:"elapsed_time"`         // seconds
	TotalElevation   float64           `json:"total_elevation_gain"` // meters
	AverageSpeed     float64           `json:"average_speed"`        // m/s
	MaxSpeed         float64           `json:"max_speed"`            // m/s
	AverageHeartrate float64           `json:"average_heartrate"`
	MaxHeartrate     float64           `json:"max_heartrate"`
	Calories         float64           `json:"calories"`
	StartDate        string            `json:"start_date"` // ISO 8601
	StartLatlng      []float64         `json:"start_latlng"`
	EndLatlng        []float64         `json:"end_latlng"`
	Map              stravaActivityMap `json:"map"`
}

type stravaActivityMap struct {
	SummaryPolyline string `json:"summary_polyline"`
}

// ── OAuth ──

func (s *stravaService) GetAuthURL(userId uuid.UUID) (response.StravaAuthURLResponse, error) {
	clientId := os.Getenv("STRAVA_CLIENT_ID")
	redirectURI := os.Getenv("STRAVA_REDIRECT_URI")
	if clientId == "" || redirectURI == "" {
		return response.StravaAuthURLResponse{}, errors.New("STRAVA_CLIENT_ID atau STRAVA_REDIRECT_URI belum dikonfigurasi")
	}

	authURL := fmt.Sprintf(
		"https://www.strava.com/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s&approval_prompt=auto&scope=read,activity:read_all&state=%s",
		clientId,
		url.QueryEscape(redirectURI),
		userId.String(),
	)

	return response.StravaAuthURLResponse{AuthURL: authURL}, nil
}

func (s *stravaService) HandleCallback(userId uuid.UUID, code string) (response.StravaConnectionResponse, error) {
	clientId := os.Getenv("STRAVA_CLIENT_ID")
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")
	if clientId == "" || clientSecret == "" {
		return response.StravaConnectionResponse{}, errors.New("Strava credentials belum dikonfigurasi")
	}

	// Exchange code for tokens
	tokenRes, err := exchangeToken(clientId, clientSecret, code)
	if err != nil {
		return response.StravaConnectionResponse{}, fmt.Errorf("gagal exchange token: %w", err)
	}

	// Check if user already has a connection
	existing, _ := s.stravaRepo.FindConnectionByUserId(userId)
	if existing != nil && existing.Id != uuid.Nil {
		// Update existing connection
		existing.AthleteId = tokenRes.Athlete.Id
		existing.AccessToken = tokenRes.AccessToken
		existing.RefreshToken = tokenRes.RefreshToken
		existing.ExpiresAt = tokenRes.ExpiresAt
		existing.Scope = "read,activity:read_all"
		existing.IsConnected = true
		existing.UpdatedAt = time.Now()

		if err := s.stravaRepo.UpdateConnection(existing); err != nil {
			return response.StravaConnectionResponse{}, err
		}
		return toConnectionResponse(*existing), nil
	}

	// Create new connection
	conn := entity.StravaConnection{
		Id:           uuid.New(),
		UserId:       userId,
		AthleteId:    tokenRes.Athlete.Id,
		AccessToken:  tokenRes.AccessToken,
		RefreshToken: tokenRes.RefreshToken,
		ExpiresAt:    tokenRes.ExpiresAt,
		Scope:        "read,activity:read_all",
		IsConnected:  true,
		ConnectedAt:  time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.stravaRepo.CreateConnection(&conn); err != nil {
		return response.StravaConnectionResponse{}, err
	}

	return toConnectionResponse(conn), nil
}

func (s *stravaService) Disconnect(userId uuid.UUID) error {
	conn, err := s.stravaRepo.FindConnectionByUserId(userId)
	if err != nil {
		return errors.New("koneksi Strava tidak ditemukan")
	}

	// Deauthorize on Strava side
	_ = deauthorizeStrava(conn.AccessToken)

	conn.IsConnected = false
	conn.AccessToken = ""
	conn.RefreshToken = ""
	conn.UpdatedAt = time.Now()
	return s.stravaRepo.UpdateConnection(conn)
}

func (s *stravaService) GetConnection(userId uuid.UUID) (response.StravaConnectionResponse, error) {
	conn, err := s.stravaRepo.FindConnectionByUserId(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.StravaConnectionResponse{}, errors.New("belum terhubung dengan Strava")
		}
		return response.StravaConnectionResponse{}, err
	}
	return toConnectionResponse(*conn), nil
}

// ── Activity Sync ──

func (s *stravaService) SyncActivities(userId uuid.UUID) (response.StravaSyncSummaryResponse, error) {
	conn, err := s.stravaRepo.FindConnectionByUserId(userId)
	if err != nil {
		return response.StravaSyncSummaryResponse{}, errors.New("belum terhubung dengan Strava")
	}
	if !conn.IsConnected {
		return response.StravaSyncSummaryResponse{}, errors.New("koneksi Strava tidak aktif")
	}

	// Refresh token if expired
	if err := s.refreshTokenIfNeeded(conn); err != nil {
		return response.StravaSyncSummaryResponse{}, fmt.Errorf("gagal refresh token: %w", err)
	}

	// Fetch activities from Strava (last 30 days, max 50)
	after := time.Now().AddDate(0, 0, -30).Unix()
	activities, err := fetchStravaActivities(conn.AccessToken, after, 50)
	if err != nil {
		return response.StravaSyncSummaryResponse{}, fmt.Errorf("gagal mengambil aktivitas dari Strava: %w", err)
	}

	summary := response.StravaSyncSummaryResponse{}

	for _, act := range activities {
		// Only sync running activities
		if !isRunningActivity(act.Type, act.SportType) {
			continue
		}

		// Check for duplicate
		existing, _ := s.stravaRepo.FindActivityByStravaId(act.Id)
		if existing != nil && existing.Id != uuid.Nil {
			summary.TotalSkipped++
			continue
		}

		// Parse start date
		startDate, _ := time.Parse(time.RFC3339, act.StartDate)
		if startDate.IsZero() {
			startDate, _ = time.Parse("2006-01-02T15:04:05Z", act.StartDate)
		}

		// Convert to km + pace
		distanceKm := act.Distance / 1000.0
		var avgPace float64
		if distanceKm > 0 && act.MovingTime > 0 {
			avgPace = (float64(act.MovingTime) / 60.0) / distanceKm // min/km
			avgPace = math.Round(avgPace*100) / 100
		}

		// Create linked RunActivity
		runActivity := entity.RunActivity{
			Id:        uuid.New(),
			UserId:    userId,
			Distance:  distanceKm,
			Duration:  act.MovingTime,
			AvgPace:   avgPace,
			Calories:  int(act.Calories),
			Source:    "strava",
			CreatedAt: time.Now(),
		}
		if err := s.activityRepo.Create(&runActivity); err != nil {
			summary.TotalFailed++
			continue
		}

		// Extract coordinates
		var startLat, startLng, endLat, endLng float64
		if len(act.StartLatlng) == 2 {
			startLat = act.StartLatlng[0]
			startLng = act.StartLatlng[1]
		}
		if len(act.EndLatlng) == 2 {
			endLat = act.EndLatlng[0]
			endLng = act.EndLatlng[1]
		}

		stravaActivity := entity.StravaActivity{
			Id:               uuid.New(),
			ConnectionId:     conn.Id,
			UserId:           userId,
			RunActivityId:    runActivity.Id,
			StravaId:         act.Id,
			Name:             act.Name,
			Type:             act.Type,
			Distance:         act.Distance,
			MovingTime:       act.MovingTime,
			ElapsedTime:      act.ElapsedTime,
			TotalElevation:   act.TotalElevation,
			AverageSpeed:     act.AverageSpeed,
			MaxSpeed:         act.MaxSpeed,
			AverageHeartrate: act.AverageHeartrate,
			MaxHeartrate:     act.MaxHeartrate,
			Calories:         act.Calories,
			StartDate:        startDate,
			StartLatitude:    startLat,
			StartLongitude:   startLng,
			EndLatitude:      endLat,
			EndLongitude:     endLng,
			MapPolyline:      act.Map.SummaryPolyline,
			Status:           "synced",
			SyncedAt:         time.Now(),
			CreatedAt:        time.Now(),
		}

		if err := s.stravaRepo.CreateActivity(&stravaActivity); err != nil {
			summary.TotalFailed++
			continue
		}

		summary.TotalSynced++
		summary.Activities = append(summary.Activities, toActivityResponse(stravaActivity))
	}

	// Update last sync time
	now := time.Now()
	conn.LastSyncAt = &now
	_ = s.stravaRepo.UpdateConnection(conn)

	return summary, nil
}

func (s *stravaService) GetSyncHistory(userId uuid.UUID, limit int) ([]response.StravaActivityResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	activities, err := s.stravaRepo.FindActivitiesByUserId(userId, limit)
	if err != nil {
		return nil, err
	}

	var res []response.StravaActivityResponse
	for _, a := range activities {
		res = append(res, toActivityResponse(a))
	}
	return res, nil
}

func (s *stravaService) GetStats(userId uuid.UUID) (response.StravaStatsResponse, error) {
	conn, err := s.stravaRepo.FindConnectionByUserId(userId)
	if err != nil {
		return response.StravaStatsResponse{}, errors.New("belum terhubung dengan Strava")
	}

	stats, err := s.stravaRepo.GetSyncStats(userId)
	if err != nil {
		return response.StravaStatsResponse{}, err
	}

	result := response.StravaStatsResponse{
		Connection: toConnectionResponse(*conn),
	}

	if stats != nil {
		result.TotalActivities = toInt(stats["total_activities"])
		totalDistanceM := toFloat(stats["total_distance"])
		result.TotalDistanceKm = math.Round((totalDistanceM/1000.0)*100) / 100
		result.TotalDuration = toInt(stats["total_duration"])

		avgSpeed := toFloat(stats["avg_speed"])
		if avgSpeed > 0 {
			// Convert m/s to min/km
			result.AvgPace = math.Round(((1000.0/avgSpeed)/60.0)*100) / 100
		}

		result.AvgHeartrate = math.Round(toFloat(stats["avg_heartrate"])*10) / 10
		result.TotalCalories = toInt(stats["total_calories"])
		result.TotalElevation = math.Round(toFloat(stats["total_elevation"])*100) / 100
	}

	return result, nil
}

// ── Token management ──

func (s *stravaService) refreshTokenIfNeeded(conn *entity.StravaConnection) error {
	if time.Now().Unix() < conn.ExpiresAt-60 {
		return nil // Still valid
	}

	clientId := os.Getenv("STRAVA_CLIENT_ID")
	clientSecret := os.Getenv("STRAVA_CLIENT_SECRET")

	data := url.Values{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {conn.RefreshToken},
	}

	resp, err := http.PostForm("https://www.strava.com/oauth/token", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("strava refresh token failed: %s", string(body))
	}

	var tokenRes stravaTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return err
	}

	conn.AccessToken = tokenRes.AccessToken
	conn.RefreshToken = tokenRes.RefreshToken
	conn.ExpiresAt = tokenRes.ExpiresAt
	conn.UpdatedAt = time.Now()

	return s.stravaRepo.UpdateConnection(conn)
}

// ── HTTP helpers ──

func exchangeToken(clientId, clientSecret, code string) (*stravaTokenResponse, error) {
	data := url.Values{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	}

	resp, err := http.PostForm("https://www.strava.com/oauth/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("strava token exchange failed (%d): %s", resp.StatusCode, string(body))
	}

	var tokenRes stravaTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return nil, err
	}

	return &tokenRes, nil
}

func fetchStravaActivities(accessToken string, after int64, perPage int) ([]stravaActivityAPI, error) {
	apiURL := fmt.Sprintf(
		"https://www.strava.com/api/v3/athlete/activities?after=%d&per_page=%d",
		after, perPage,
	)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("strava API error (%d): %s", resp.StatusCode, string(body))
	}

	var activities []stravaActivityAPI
	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		return nil, err
	}

	return activities, nil
}

func deauthorizeStrava(accessToken string) error {
	data := url.Values{"access_token": {accessToken}}
	resp, err := http.PostForm("https://www.strava.com/oauth/deauthorize", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func isRunningActivity(actType, sportType string) bool {
	runTypes := []string{"Run", "VirtualRun", "TrailRun"}
	for _, t := range runTypes {
		if strings.EqualFold(actType, t) || strings.EqualFold(sportType, t) {
			return true
		}
	}
	return false
}

// ── Response mappers ──

func toConnectionResponse(c entity.StravaConnection) response.StravaConnectionResponse {
	return response.StravaConnectionResponse{
		Id:          c.Id.String(),
		AthleteId:   c.AthleteId,
		Scope:       c.Scope,
		IsConnected: c.IsConnected,
		LastSyncAt:  c.LastSyncAt,
		ConnectedAt: c.ConnectedAt,
	}
}

func toActivityResponse(a entity.StravaActivity) response.StravaActivityResponse {
	distanceKm := math.Round((a.Distance/1000.0)*100) / 100

	var avgPace float64
	if a.AverageSpeed > 0 {
		avgPace = math.Round(((1000.0/a.AverageSpeed)/60.0)*100) / 100
	}

	return response.StravaActivityResponse{
		Id:               a.Id.String(),
		StravaId:         a.StravaId,
		RunActivityId:    a.RunActivityId.String(),
		Name:             a.Name,
		Type:             a.Type,
		DistanceKm:       distanceKm,
		MovingTime:       a.MovingTime,
		ElapsedTime:      a.ElapsedTime,
		TotalElevation:   a.TotalElevation,
		AveragePace:      avgPace,
		MaxSpeed:         a.MaxSpeed,
		AverageHeartrate: a.AverageHeartrate,
		MaxHeartrate:     a.MaxHeartrate,
		Calories:         a.Calories,
		StartDate:        a.StartDate,
		MapPolyline:      a.MapPolyline,
		Status:           a.Status,
		SyncedAt:         a.SyncedAt,
	}
}

func toInt(v interface{}) int {
	switch n := v.(type) {
	case int64:
		return int(n)
	case float64:
		return int(n)
	case int:
		return n
	default:
		return 0
	}
}

func toFloat(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int64:
		return float64(n)
	case int:
		return float64(n)
	default:
		return 0
	}
}
