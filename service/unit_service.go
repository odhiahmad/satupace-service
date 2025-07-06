package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type UnitService interface {
	Create(req request.UnitRequest) (entity.Unit, error)
	Update(id int, req request.UnitRequest) (entity.Unit, error)
	Delete(id int) error
	FindById(id int) response.UnitResponse
	FindWithPagination(businessId int, pagination request.Pagination) ([]response.UnitResponse, int64, error)
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

func (s *unitService) Create(req request.UnitRequest) (entity.Unit, error) {
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

func (s *unitService) Update(id int, req request.UnitRequest) (entity.Unit, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Unit{}, err
	}

	unit := entity.Unit{
		Id:         id,
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

func (s *unitService) FindById(unitId int) response.UnitResponse {
	unitData, err := s.repo.FindById(unitId)
	helper.ErrorPanic(err)

	unitResponse := helper.MapUnit(&unitData)
	return *unitResponse
}

func (s *unitService) FindWithPagination(businessId int, pagination request.Pagination) ([]response.UnitResponse, int64, error) {
	units, total, err := s.repo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.UnitResponse
	for _, unit := range units {
		result = append(result, *helper.MapUnit(&unit))
	}

	return result, total, nil
}
