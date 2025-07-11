package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/repository"
	"github.com/redis/go-redis/v9"
)

type LocationService interface {
	GetProvinces() ([]response.ProvinceResponse, error)
	GetCitiesByProvince(provinceID int) ([]response.CityResponse, error)
	GetDistrictsByCity(cityID int) ([]response.DistrictResponse, error)
	GetVillagesByDistrict(districtID int) ([]response.VillageResponse, error)
}

type locationService struct {
	locationRepo repository.LocationRepository
	redisClient  *redis.Client
}

func NewLocationService(repo repository.LocationRepository, redisClient *redis.Client) LocationService {
	return &locationService{locationRepo: repo, redisClient: redisClient}
}

var ctx = context.Background()

func (s *locationService) GetProvinces() ([]response.ProvinceResponse, error) {
	cacheKey := "provinces"
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var res []response.ProvinceResponse
		if err := json.Unmarshal([]byte(cached), &res); err == nil {
			return res, nil
		}
	}

	provinces, err := s.locationRepo.GetProvinces()
	if err != nil {
		return nil, err
	}

	var res []response.ProvinceResponse
	for _, p := range provinces {
		res = append(res, response.ProvinceResponse{
			ID:   p.ID,
			Name: p.Name,
			Code: p.Code,
		})
	}

	data, _ := json.Marshal(res)
	s.redisClient.Set(ctx, cacheKey, data, 6*time.Hour)

	return res, nil
}

func (s *locationService) GetCitiesByProvince(provinceID int) ([]response.CityResponse, error) {
	cacheKey := fmt.Sprintf("cities:province:%d", provinceID)
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var res []response.CityResponse
		if err := json.Unmarshal([]byte(cached), &res); err == nil {
			return res, nil
		}
	}

	cities, err := s.locationRepo.GetCitiesByProvinceID(provinceID)
	if err != nil {
		return nil, err
	}

	var res []response.CityResponse
	for _, c := range cities {
		res = append(res, response.CityResponse{
			ID:         c.ID,
			ProvinceID: c.ProvinceID,
			Type:       c.Type,
			Name:       c.Name,
			Code:       c.Code,
			FullCode:   c.FullCode,
		})
	}

	data, _ := json.Marshal(res)
	s.redisClient.Set(ctx, cacheKey, data, 6*time.Hour)

	return res, nil
}

func (s *locationService) GetDistrictsByCity(cityID int) ([]response.DistrictResponse, error) {
	cacheKey := fmt.Sprintf("districts:city:%d", cityID)
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var res []response.DistrictResponse
		if err := json.Unmarshal([]byte(cached), &res); err == nil {
			return res, nil
		}
	}

	districts, err := s.locationRepo.GetDistrictsByCityID(cityID)
	if err != nil {
		return nil, err
	}

	var res []response.DistrictResponse
	for _, d := range districts {
		res = append(res, response.DistrictResponse{
			ID:       d.ID,
			CityID:   d.CityID,
			Name:     d.Name,
			Code:     d.Code,
			FullCode: d.FullCode,
		})
	}

	data, _ := json.Marshal(res)
	s.redisClient.Set(ctx, cacheKey, data, 6*time.Hour)

	return res, nil
}

func (s *locationService) GetVillagesByDistrict(districtID int) ([]response.VillageResponse, error) {
	cacheKey := fmt.Sprintf("villages:district:%d", districtID)
	cached, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var res []response.VillageResponse
		if err := json.Unmarshal([]byte(cached), &res); err == nil {
			return res, nil
		}
	}

	villages, err := s.locationRepo.GetVillagesByDistrictID(districtID)
	if err != nil {
		return nil, err
	}

	var res []response.VillageResponse
	for _, v := range villages {
		res = append(res, response.VillageResponse{
			ID:         v.ID,
			DistrictID: v.DistrictID,
			Name:       v.Name,
			Code:       v.Code,
			FullCode:   v.FullCode,
			PosCode:    v.PosCode,
		})
	}

	data, _ := json.Marshal(res)
	s.redisClient.Set(ctx, cacheKey, data, 6*time.Hour)

	return res, nil
}
