package service

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type DiscountService interface {
	Create(req request.DiscountCreate) (entity.Discount, error)
	Update(id int, req request.DiscountUpdate) (entity.Discount, error)
	Delete(id int) error
	FindById(id int) (response.DiscountResponse, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.DiscountResponse, int64, error)
	SetIsActive(id int, active bool) error
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

func (s *discountService) Create(req request.DiscountCreate) (entity.Discount, error) {
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
		Name:         req.Name,
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

func (s *discountService) Update(id int, req request.DiscountUpdate) (entity.Discount, error) {
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
		Name:         req.Name,
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

func (s *discountService) FindById(id int) (response.DiscountResponse, error) {
	discount, err := s.repo.FindById(id)
	if err != nil {
		return response.DiscountResponse{}, err
	}
	return ToDiscountResponse(discount), nil
}

func (s *discountService) FindWithPagination(businessId int, pagination request.Pagination) ([]response.DiscountResponse, int64, error) {
	discounts, total, err := s.repo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var responses []response.DiscountResponse
	for _, d := range discounts {
		responses = append(responses, ToDiscountResponse(d)) // fungsi mapping
	}

	return responses, total, nil
}

func (s *discountService) SetIsActive(id int, active bool) error {
	return s.repo.SetIsActive(id, active)
}

func ToDiscountResponse(discount entity.Discount) response.DiscountResponse {
	return response.DiscountResponse{
		Id:           discount.Id,
		Name:         discount.Name,
		Description:  discount.Description,
		IsPercentage: discount.IsPercentage,
		Amount:       discount.Amount,
		IsGlobal:     discount.IsGlobal,
		IsMultiple:   discount.IsMultiple,
		StartAt:      discount.StartAt,
		EndAt:        discount.EndAt,
		IsActive:     discount.IsActive,
	}
}
