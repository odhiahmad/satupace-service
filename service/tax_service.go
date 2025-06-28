package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type TaxService interface {
	Create(req request.TaxCreate) (entity.Tax, error)
	Update(id int, req request.TaxUpdate) (entity.Tax, error)
	Delete(id int) error
	FindById(id int) (entity.Tax, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Tax, int64, error)
}

type taxService struct {
	repo     repository.TaxRepository
	validate *validator.Validate
}

func NewTaxService(repo repository.TaxRepository, validate *validator.Validate) TaxService {
	return &taxService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *taxService) Create(req request.TaxCreate) (entity.Tax, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Tax{}, err
	}

	typeVal := helper.DeterminePromoType(req.Amount)

	tax := entity.Tax{
		BusinessId: req.BusinessId,
		Name:       req.Name,
		Type:       typeVal,
		Amount:     req.Amount,
		IsGlobal:   req.IsGlobal,
	}

	return s.repo.Create(tax)
}

func (s *taxService) Update(id int, req request.TaxUpdate) (entity.Tax, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Tax{}, err
	}

	typeVal := helper.DeterminePromoType(req.Amount)

	tax := entity.Tax{
		Id:       id,
		Name:     req.Name,
		Type:     typeVal,
		Amount:   req.Amount,
		IsGlobal: req.IsGlobal,
	}

	return s.repo.Update(tax)
}

func (s *taxService) Delete(id int) error {
	tax, err := s.repo.FindById(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(tax)
}

func (s *taxService) FindById(id int) (entity.Tax, error) {
	return s.repo.FindById(id)
}

func (s *taxService) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Tax, int64, error) {
	return s.repo.FindWithPagination(businessId, pagination)
}
