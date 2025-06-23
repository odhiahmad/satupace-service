package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/repository"
)

type ProductUnitService interface {
	Create(req request.ProductUnitCreate) (entity.ProductUnit, error)
	Update(req request.ProductUnitUpdate) (entity.ProductUnit, error)
	Delete(id int) error
	FindById(id int) (entity.ProductUnit, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.ProductUnit, int64, error)
}

type productUnitService struct {
	repo     repository.ProductUnitRepository
	validate *validator.Validate
}

func NewProductUnitService(repo repository.ProductUnitRepository) ProductUnitService {
	return &productUnitService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *productUnitService) Create(req request.ProductUnitCreate) (entity.ProductUnit, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.ProductUnit{}, err
	}

	productUnit := entity.ProductUnit{
		BusinessId: req.BusinessId,
		Name:       req.Name,
		Alias:      req.Alias,
		Multiplier: req.Multiplier,
	}

	return s.repo.Create(productUnit)
}

func (s *productUnitService) Update(req request.ProductUnitUpdate) (entity.ProductUnit, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.ProductUnit{}, err
	}

	productUnit := entity.ProductUnit{
		Id:         req.Id,
		BusinessId: req.BusinessId,
		Name:       req.Name,
		Alias:      req.Alias,
		Multiplier: req.Multiplier,
	}

	return s.repo.Update(productUnit)
}

func (s *productUnitService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *productUnitService) FindById(id int) (entity.ProductUnit, error) {
	return s.repo.FindById(id)
}

func (s *productUnitService) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.ProductUnit, int64, error) {
	return s.repo.FindWithPagination(businessId, pagination)
}
