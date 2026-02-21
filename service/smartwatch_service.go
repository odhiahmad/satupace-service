package service

import (
	"errors"
	"run-sync/data/request"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SmartWatchService interface {
	ConnectDevice(userId uuid.UUID, req request.ConnectDeviceRequest) (response.SmartWatchDeviceResponse, error)
	DisconnectDevice(userId uuid.UUID, deviceId uuid.UUID) error
	GetDevices(userId uuid.UUID) ([]response.SmartWatchDeviceResponse, error)
	SyncActivity(userId uuid.UUID, req request.SyncActivityRequest) (response.SmartWatchSyncResponse, error)
	BatchSync(userId uuid.UUID, req request.BatchSyncRequest) (response.SyncSummaryResponse, error)
	GetSyncHistory(userId uuid.UUID, limit int) ([]response.SmartWatchSyncResponse, error)
	GetDeviceStats(userId uuid.UUID, deviceId uuid.UUID) (response.DeviceStatsResponse, error)
}

type smartWatchService struct {
	watchRepo    repository.SmartWatchRepository
	activityRepo repository.RunActivityRepository
}

func NewSmartWatchService(watchRepo repository.SmartWatchRepository, activityRepo repository.RunActivityRepository) SmartWatchService {
	return &smartWatchService{watchRepo: watchRepo, activityRepo: activityRepo}
}

func (s *smartWatchService) ConnectDevice(userId uuid.UUID, req request.ConnectDeviceRequest) (response.SmartWatchDeviceResponse, error) {
	device := entity.SmartWatchDevice{
		Id:           uuid.New(),
		UserId:       userId,
		DeviceType:   req.DeviceType,
		DeviceName:   req.DeviceName,
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
		ExternalId:   req.ExternalId,
		IsConnected:  true,
		ConnectedAt:  time.Now(),
	}

	if err := s.watchRepo.CreateDevice(&device); err != nil {
		return response.SmartWatchDeviceResponse{}, err
	}

	return toDeviceResponse(device), nil
}

func (s *smartWatchService) DisconnectDevice(userId uuid.UUID, deviceId uuid.UUID) error {
	device, err := s.watchRepo.FindDeviceById(deviceId)
	if err != nil {
		return err
	}
	if device.UserId != userId {
		return errors.New("unauthorized: device does not belong to user")
	}

	device.IsConnected = false
	device.AccessToken = ""
	device.RefreshToken = ""
	return s.watchRepo.UpdateDevice(device)
}

func (s *smartWatchService) GetDevices(userId uuid.UUID) ([]response.SmartWatchDeviceResponse, error) {
	devices, err := s.watchRepo.FindDevicesByUserId(userId)
	if err != nil {
		return nil, err
	}

	var res []response.SmartWatchDeviceResponse
	for _, d := range devices {
		res = append(res, toDeviceResponse(d))
	}
	return res, nil
}

func (s *smartWatchService) SyncActivity(userId uuid.UUID, req request.SyncActivityRequest) (response.SmartWatchSyncResponse, error) {
	// Parse device ID
	deviceId, err := uuid.Parse(req.DeviceId)
	if err != nil {
		return response.SmartWatchSyncResponse{}, errors.New("invalid device_id")
	}

	// Verify device ownership
	device, err := s.watchRepo.FindDeviceById(deviceId)
	if err != nil {
		return response.SmartWatchSyncResponse{}, errors.New("device not found")
	}
	if device.UserId != userId {
		return response.SmartWatchSyncResponse{}, errors.New("unauthorized: device does not belong to user")
	}

	// Check for duplicate external_id
	existing, err := s.watchRepo.FindSyncByExternalId(req.ExternalId)
	if err == nil && existing != nil {
		return toSyncResponse(*existing), nil // Idempotent: return existing
	}

	// Parse activity date
	activityDate, err := time.Parse(time.RFC3339, req.ActivityDate)
	if err != nil {
		return response.SmartWatchSyncResponse{}, errors.New("invalid activity_date: must be RFC3339 format")
	}

	// Create a linked RunActivity
	activity := entity.RunActivity{
		Id:        uuid.New(),
		UserId:    userId,
		Distance:  req.Distance,
		Duration:  req.Duration,
		AvgPace:   req.AvgPace,
		Calories:  req.Calories,
		Source:    "smartwatch:" + device.DeviceType,
		CreatedAt: time.Now(),
	}
	if err := s.activityRepo.Create(&activity); err != nil {
		return response.SmartWatchSyncResponse{}, err
	}

	// Create sync record
	syncRecord := entity.SmartWatchSync{
		Id:             uuid.New(),
		DeviceId:       deviceId,
		UserId:         userId,
		ActivityId:     activity.Id,
		ExternalId:     req.ExternalId,
		RawData:        req.RawData,
		Distance:       req.Distance,
		Duration:       req.Duration,
		AvgPace:        req.AvgPace,
		MaxPace:        req.MaxPace,
		AvgHeartRate:   req.AvgHeartRate,
		MaxHeartRate:   req.MaxHeartRate,
		Calories:       req.Calories,
		Cadence:        req.Cadence,
		ElevationGain:  req.ElevationGain,
		StartLatitude:  req.StartLatitude,
		StartLongitude: req.StartLongitude,
		EndLatitude:    req.EndLatitude,
		EndLongitude:   req.EndLongitude,
		RouteData:      req.RouteData,
		Status:         "synced",
		SyncedAt:       time.Now(),
		ActivityDate:   activityDate,
		CreatedAt:      time.Now(),
	}

	if err := s.watchRepo.CreateSync(&syncRecord); err != nil {
		return response.SmartWatchSyncResponse{}, err
	}

	// Update device last sync time
	now := time.Now()
	device.LastSyncAt = &now
	_ = s.watchRepo.UpdateDevice(device)

	return toSyncResponse(syncRecord), nil
}

func (s *smartWatchService) BatchSync(userId uuid.UUID, req request.BatchSyncRequest) (response.SyncSummaryResponse, error) {
	summary := response.SyncSummaryResponse{}

	for _, actReq := range req.Activities {
		actReq.DeviceId = req.DeviceId // Ensure device ID from batch wrapper
		res, err := s.SyncActivity(userId, actReq)
		if err != nil {
			summary.TotalFailed++
			continue
		}
		summary.TotalSynced++
		summary.Activities = append(summary.Activities, res)
	}

	return summary, nil
}

func (s *smartWatchService) GetSyncHistory(userId uuid.UUID, limit int) ([]response.SmartWatchSyncResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	syncs, err := s.watchRepo.FindSyncsByUserId(userId, limit)
	if err != nil {
		return nil, err
	}

	var res []response.SmartWatchSyncResponse
	for _, sync := range syncs {
		res = append(res, toSyncResponse(sync))
	}
	return res, nil
}

func (s *smartWatchService) GetDeviceStats(userId uuid.UUID, deviceId uuid.UUID) (response.DeviceStatsResponse, error) {
	device, err := s.watchRepo.FindDeviceById(deviceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.DeviceStatsResponse{}, errors.New("device not found")
		}
		return response.DeviceStatsResponse{}, err
	}
	if device.UserId != userId {
		return response.DeviceStatsResponse{}, errors.New("unauthorized: device does not belong to user")
	}

	stats, err := s.watchRepo.GetDeviceStats(deviceId)
	if err != nil {
		return response.DeviceStatsResponse{}, err
	}

	result := response.DeviceStatsResponse{
		Device: toDeviceResponse(*device),
	}

	if stats != nil {
		if v, ok := stats["total_activities"]; ok {
			result.TotalActivities = toInt(v)
		}
		if v, ok := stats["total_distance"]; ok {
			result.TotalDistance = toFloat(v)
		}
		if v, ok := stats["total_duration"]; ok {
			result.TotalDuration = toInt(v)
		}
		if v, ok := stats["avg_pace"]; ok {
			result.AvgPace = toFloat(v)
		}
		if v, ok := stats["avg_heart_rate"]; ok {
			result.AvgHeartRate = toFloat(v)
		}
		if v, ok := stats["total_calories"]; ok {
			result.TotalCalories = toInt(v)
		}
		if v, ok := stats["total_elevation"]; ok {
			result.TotalElevation = toFloat(v)
		}
	}

	return result, nil
}

// ── Helpers ──

func toDeviceResponse(d entity.SmartWatchDevice) response.SmartWatchDeviceResponse {
	return response.SmartWatchDeviceResponse{
		Id:          d.Id.String(),
		UserId:      d.UserId.String(),
		DeviceType:  d.DeviceType,
		DeviceName:  d.DeviceName,
		ExternalId:  d.ExternalId,
		IsConnected: d.IsConnected,
		LastSyncAt:  d.LastSyncAt,
		ConnectedAt: d.ConnectedAt,
	}
}

func toSyncResponse(s entity.SmartWatchSync) response.SmartWatchSyncResponse {
	return response.SmartWatchSyncResponse{
		Id:             s.Id.String(),
		DeviceId:       s.DeviceId.String(),
		UserId:         s.UserId.String(),
		ActivityId:     s.ActivityId.String(),
		ExternalId:     s.ExternalId,
		Distance:       s.Distance,
		Duration:       s.Duration,
		AvgPace:        s.AvgPace,
		MaxPace:        s.MaxPace,
		AvgHeartRate:   s.AvgHeartRate,
		MaxHeartRate:   s.MaxHeartRate,
		Calories:       s.Calories,
		Cadence:        s.Cadence,
		ElevationGain:  s.ElevationGain,
		StartLatitude:  s.StartLatitude,
		StartLongitude: s.StartLongitude,
		EndLatitude:    s.EndLatitude,
		EndLongitude:   s.EndLongitude,
		RouteData:      s.RouteData,
		Status:         s.Status,
		ActivityDate:   s.ActivityDate,
		SyncedAt:       s.SyncedAt,
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
