package service

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type DiscountService interface {
	Create(req request.DiscountRequest) (entity.Discount, error)
	Update(id int, req request.DiscountRequest) (entity.Discount, error)
	Delete(id int) error
	SetIsActive(id int, active bool) error
	FindById(id int) (response.DiscountResponse, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.DiscountResponse, int64, error)
	FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]response.DiscountResponse, string, bool, error)
}

type discountService struct {
	repo     repository.DiscountRepository
	validate *validator.Validate
}

func NewDiscountService(repo repository.DiscountRepository, validate *validator.Validate) DiscountService {
	return &discountService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *discountService) Create(req request.DiscountRequest) (entity.Discount, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Discount{}, err
	}
	if req.Amount <= 0 {
		return entity.Discount{}, errors.New("jumlah diskon harus lebih besar dari 0")
	}
	if !req.EndAt.After(req.StartAt) {
		return entity.Discount{}, errors.New("tanggal berakhir harus setelah tanggal mulai")
	}

	isPercentageVal := helper.DeterminePromoType(req.Amount)

	discount := entity.Discount{
		BusinessId:   req.BusinessId,
		Name:         strings.ToLower(req.Name),
		Description:  req.Description,
		IsPercentage: isPercentageVal,
		Amount:       req.Amount,
		StartAt:      req.StartAt,
		EndAt:        req.EndAt,
		IsGlobal:     req.IsGlobal,
		IsMultiple:   req.IsMultiple,
		IsActive:     true,
	}

	return s.repo.Create(discount)
}

func (s *discountService) Update(id int, req request.DiscountRequest) (entity.Discount, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Discount{}, err
	}
	if req.Amount <= 0 {
		return entity.Discount{}, errors.New("jumlah diskon harus lebih besar dari 0")
	}
	if !req.EndAt.After(req.StartAt) {
		return entity.Discount{}, errors.New("tanggal berakhir harus setelah tanggal mulai")
	}

	isPercentageVal := helper.DeterminePromoType(req.Amount)

	discount := entity.Discount{
		Id:           id,
		Name:         strings.ToLower(req.Name),
		Description:  req.Description,
		IsPercentage: isPercentageVal,
		Amount:       req.Amount,
		StartAt:      req.StartAt,
		EndAt:        req.EndAt,
		IsGlobal:     req.IsGlobal,
		IsMultiple:   req.IsMultiple,
	}

	return s.repo.Update(discount)
}

func (s *discountService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *discountService) SetIsActive(id int, active bool) error {
	return s.repo.SetIsActive(id, active)
}

func (s *discountService) FindById(id int) (response.DiscountResponse, error) {
	discount, err := s.repo.FindById(id)
	if err != nil {
		return response.DiscountResponse{}, err
	}
	return helper.ToDiscountResponse(discount), nil
}

func (s *discountService) FindWithPagination(businessId int, pagination request.Pagination) ([]response.DiscountResponse, int64, error) {
	discounts, total, err := s.repo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var responses []response.DiscountResponse
	for _, d := range discounts {
		responses = append(responses, helper.ToDiscountResponse(d)) // fungsi mapping
	}

	return responses, total, nil
}

func (s *discountService) FindWithPaginationCursor(businessId int, pagination request.Pagination) ([]response.DiscountResponse, string, bool, error) {
	discounts, nextCursor, hasNext, err := s.repo.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var responses []response.DiscountResponse
	for _, d := range discounts {
		responses = append(responses, helper.ToDiscountResponse(d))
	}

	return responses, nextCursor, hasNext, nil
}
