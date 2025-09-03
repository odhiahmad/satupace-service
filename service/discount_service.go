package service

import (
	"errors"
	"strings"

	"loka-kasir/data/request"
	"loka-kasir/data/response"
	"loka-kasir/entity"
	"loka-kasir/helper"
	"loka-kasir/helper/mapper"
	"loka-kasir/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type DiscountService interface {
	Create(req request.DiscountRequest) (entity.Discount, error)
	Update(id uuid.UUID, req request.DiscountRequest) (entity.Discount, error)
	Delete(id uuid.UUID) error
	SetIsActive(id uuid.UUID, active bool) error
	FindById(id uuid.UUID) (response.DiscountResponse, error)
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.DiscountResponse, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.DiscountResponse, string, bool, error)
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
		Description:  strings.ToLower(req.Description),
		IsPercentage: &isPercentageVal,
		Amount:       req.Amount,
		StartAt:      req.StartAt,
		EndAt:        req.EndAt,
		IsGlobal:     req.IsGlobal,
		IsMultiple:   req.IsMultiple,
		IsActive:     req.IsActive,
	}

	return s.repo.Create(discount)
}

func (s *discountService) Update(id uuid.UUID, req request.DiscountRequest) (entity.Discount, error) {
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
		Description:  strings.ToLower(req.Description),
		IsPercentage: &isPercentageVal,
		Amount:       req.Amount,
		StartAt:      req.StartAt,
		EndAt:        req.EndAt,
		IsGlobal:     req.IsGlobal,
		IsMultiple:   req.IsMultiple,
		IsActive:     req.IsActive,
	}

	return s.repo.Update(discount)
}

func (s *discountService) Delete(id uuid.UUID) error {
	_, err := s.repo.FindById(id)
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

	return nil
}

func (s *discountService) SetIsActive(id uuid.UUID, active bool) error {
	return s.repo.SetIsActive(id, active)
}

func (s *discountService) FindById(id uuid.UUID) (response.DiscountResponse, error) {
	discount, err := s.repo.FindById(id)
	if err != nil {
		return response.DiscountResponse{}, err
	}
	return *mapper.MapDiscount(&discount), nil
}

func (s *discountService) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.DiscountResponse, int64, error) {
	discounts, total, err := s.repo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var responses []response.DiscountResponse
	for _, d := range discounts {
		responses = append(responses, *mapper.MapDiscount(&d)) // fungsi mapping
	}

	return responses, total, nil
}

func (s *discountService) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.DiscountResponse, string, bool, error) {
	discounts, nextCursor, hasNext, err := s.repo.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var responses []response.DiscountResponse
	for _, d := range discounts {
		responses = append(responses, *mapper.MapDiscount(&d))
	}

	return responses, nextCursor, hasNext, nil
}
