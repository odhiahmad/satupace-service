package service

import (
	"context"
	"fmt"
	"time"

	"loka-kasir/data/response"
	"loka-kasir/helper"
	"loka-kasir/helper/mapper"
	"loka-kasir/repository"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type HomeService interface {
	GetHome(businessId uuid.UUID) (*response.HomeResponse, error)
}

type homeService struct {
	repo  repository.TransactionRepository
	redis *redis.Client
}

func NewHomeService(repo repository.TransactionRepository, redis *redis.Client) HomeService {
	return &homeService{repo: repo, redis: redis}
}

func (s *homeService) GetHome(businessId uuid.UUID) (*response.HomeResponse, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("dashboard:business:%s:today", businessId.String())

	var cached response.HomeResponse
	err := helper.GetJSONFromRedis(ctx, s.redis, cacheKey, &cached)
	if err == nil {
		return &cached, nil
	}

	todayRevenue, _ := s.repo.GetTodayRevenue(businessId)
	todayOrders, _ := s.repo.GetTodayOrdersCount(businessId)
	todayItems, _ := s.repo.GetTodayItemsSold(businessId)
	lastOrders, _ := s.repo.GetLastOrders(businessId, 4)
	lastItems, _ := s.repo.GetLastItems(businessId, 5)
	topProducts, _ := s.repo.GetTopProductsToday(businessId, 5)

	var lastItemsTotal int64
	for _, i := range lastItems {
		lastItemsTotal += int64(i.Quantity)
	}

	dashboard := response.HomeResponse{
		TodaySummary: response.TodaySummaryResponse{
			TotalRevenue: todayRevenue,
			TotalOrders:  todayOrders,
			TotalItems:   todayItems,
		},
		RecentOrders:     mapper.MapTransactions(lastOrders),
		RecentItems:      mapper.MapTransactionItems(lastItems),
		RecentItemsTotal: lastItemsTotal,
		TopProducts:      topProducts,
	}

	_ = helper.SetJSONToRedis(ctx, s.redis, cacheKey, dashboard, time.Minute*5)

	return &dashboard, nil
}
