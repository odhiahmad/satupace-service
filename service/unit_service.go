package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/repository"
)

type UnitService interface {
	Create(req request.UnitCreate) (entity.Unit, error)
	Update(req request.UnitUpdate) (entity.Unit, error)
	Delete(id int) error
	FindById(id int) (entity.Unit, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Unit, int64, error)
}

type unitService struct {
	repo     repository.UnitRepository
	validate *validator.Validate
}

func NewUnitService(repo repository.UnitRepository) UnitService {
	return &unitService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *unitService) Create(req request.UnitCreate) (entity.Unit, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Unit{}, err
	}

	unit := entity.Unit{
		BusinessId: req.BusinessId,
		Name:       req.Name,
		Alias:      req.Alias,
		Multiplier: req.Multiplier,
	}

	return s.repo.Create(unit)
}

func (s *unitService) Update(req request.UnitUpdate) (entity.Unit, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Unit{}, err
	}

	unit := entity.Unit{
		Id:         req.Id,
		BusinessId: req.BusinessId,
		Name:       req.Name,
		Alias:      req.Alias,
		Multiplier: req.Multiplier,
	}

	return s.repo.Update(unit)
}

func (s *unitService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *unitService) FindById(id int) (entity.Unit, error) {
	return s.repo.FindById(id)
}

func (s *unitService) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Unit, int64, error) {
	return s.repo.FindWithPagination(businessId, pagination)
}
