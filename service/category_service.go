package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
	"github.com/redis/go-redis/v9"
)

type CategoryService interface {
	Create(category request.CategoryRequest) (response.CategoryResponse, error)
	Update(id int, category request.CategoryRequest) (response.CategoryResponse, error)
	FindById(brandId int) response.CategoryResponse
	Delete(id int) error
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.CategoryResponse, int64, error)
	FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]response.CategoryResponse, string, bool, error)
}

type categoryService struct {
	repo     repository.CategoryRepository
	Validate *validator.Validate
	Redis    *redis.Client
}

func NewCategoryService(repo repository.CategoryRepository, validate *validator.Validate, redis *redis.Client) CategoryService {
	return &categoryService{
		repo:     repo,
		Validate: validate,
		Redis:    redis,
	}
}

func (s *categoryService) Create(req request.CategoryRequest) (response.CategoryResponse, error) {
	err := s.Validate.Struct(req)
	if err != nil {
		return response.CategoryResponse{}, err
	}

	category := entity.Category{
		Name:       strings.ToLower(req.Name),
		ParentId:   req.ParentId,
		BusinessId: req.BusinessId,
		IsActive:   true,
	}

	createdCategory, err := s.repo.InsertCategory(category)
	if err != nil {
		return response.CategoryResponse{}, err
	}

	pattern := fmt.Sprintf("categories:business:%d*", req.BusinessId)
	helper.DeleteKeysByPattern(context.Background(), s.Redis, pattern)

	categoryResponse := helper.MapCategory(&createdCategory)
	return *categoryResponse, nil
}

func (s *categoryService) Update(id int, req request.CategoryRequest) (response.CategoryResponse, error) {
	err := s.Validate.Struct(req)
	if err != nil {
		return response.CategoryResponse{}, err
	}

	category, err := s.repo.FindById(id)
	if err != nil {
		return response.CategoryResponse{}, err
	}

	category.Name = strings.ToLower(req.Name)
	category.BusinessId = req.BusinessId
	category.ParentId = req.ParentId

	updatedCategory, err := s.repo.UpdateCategory(category)
	if err != nil {
		return response.CategoryResponse{}, err
	}

	cacheKey := fmt.Sprintf("category:%d", id)
	s.Redis.Del(context.Background(), cacheKey)

	categoryResponse := helper.MapCategory(&updatedCategory)
	return *categoryResponse, nil
}

func (s *categoryService) FindById(brandId int) response.CategoryResponse {
	categories, err := s.repo.FindById(brandId)
	helper.ErrorPanic(err)

	category := helper.MapCategory(&categories)
	return *category
}

func (s *categoryService) Delete(id int) error {
	category, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	hasRelation, err := s.repo.HasRelation(id)
	if err != nil {
		return err
	}

	var deleteErr error
	if hasRelation {
		deleteErr = s.repo.SoftDelete(id)
	} else {
		deleteErr = s.repo.HardDelete(id)
	}
	if deleteErr != nil {
		return deleteErr
	}

	ctx := context.Background()
	s.Redis.Del(ctx, fmt.Sprintf("category:%d", id))
	pattern := fmt.Sprintf("categories:business:%d*", category.BusinessId)
	go helper.DeleteKeysByPattern(ctx, s.Redis, pattern)

	return nil
}

func (s *categoryService) FindWithPagination(businessId int, pagination request.Pagination) ([]response.CategoryResponse, int64, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("categories:business:%d:page:%d:limit:%d:cat:%v", businessId, pagination.Page, pagination.Limit, pagination.CategoryID)

	var cached []response.CategoryResponse
	err := helper.GetJSONFromRedis(ctx, s.Redis, cacheKey, &cached)
	if err == nil {
		return cached, int64(len(cached)), nil
	}

	categories, total, err := s.repo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.CategoryResponse
	for _, category := range categories {
		result = append(result, *helper.MapCategory(&category))
	}

	_ = helper.SetJSONToRedis(ctx, s.Redis, cacheKey, result, time.Minute*10)

	return result, total, nil
}

func (s *categoryService) FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]response.CategoryResponse, string, bool, error) {
	categories, nextCursor, hasNext, err := s.repo.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var result []response.CategoryResponse
	for _, category := range categories {
		result = append(result, *helper.MapCategory(&category))
	}

	return result, nextCursor, hasNext, nil
}
