package service

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type DiscountService interface {
	Create(req request.DiscountCreate) (entity.Discount, error)
	Update(id int, req request.DiscountUpdate) (entity.Discount, error)
	Delete(id int) error
	FindById(id int) (entity.Discount, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Discount, int64, error)
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

	typeVal := helper.DeterminePromoType(req.Amount)

	discount := entity.Discount{
		BusinessId: req.BusinessId,
		Name:       req.Name,
		Type:       typeVal,
		Amount:     req.Amount,
		StartAt:    req.StartAt,
		EndAt:      req.EndAt,
		IsGlobal:   req.IsGlobal,
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

	typeVal := helper.DeterminePromoType(req.Amount)

	discount := entity.Discount{
		Id:         id,
		BusinessId: req.BusinessId,
		Name:       req.Name,
		Type:       typeVal,
		Amount:     req.Amount,
		StartAt:    req.StartAt,
		EndAt:      req.EndAt,
	}

	return s.repo.Update(discount)
}

func (s *discountService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *discountService) FindById(id int) (entity.Discount, error) {
	return s.repo.FindById(id)
}

func (s *discountService) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Discount, int64, error) {
	return s.repo.FindWithPagination(businessId, pagination)
}
