package service

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type UnitService interface {
	Create(req request.UnitRequest) (entity.Unit, error)
	Update(id uuid.UUID, req request.UnitRequest) (entity.Unit, error)
	Delete(id uuid.UUID) error
	FindById(id uuid.UUID) response.UnitResponse
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.UnitResponse, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.UnitResponse, string, bool, error)
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
		Name:       strings.ToLower(req.Name),
		Alias:      req.Alias,
		Multiplier: req.Multiplier,
	}

	return s.repo.Create(unit)
}

func (s *unitService) Update(id uuid.UUID, req request.UnitRequest) (entity.Unit, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Unit{}, err
	}

	unit := entity.Unit{
		Id:         id,
		BusinessId: req.BusinessId,
		Name:       strings.ToLower(req.Name),
		Alias:      req.Alias,
		Multiplier: req.Multiplier,
	}

	return s.repo.Update(unit)
}

func (s *unitService) Delete(id uuid.UUID) error {
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

func (s *unitService) FindById(unitId uuid.UUID) response.UnitResponse {
	unitData, err := s.repo.FindById(unitId)
	helper.ErrorPanic(err)

	unitResponse := helper.MapUnit(&unitData)
	return *unitResponse
}

func (s *unitService) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.UnitResponse, int64, error) {
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

func (s *unitService) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.UnitResponse, string, bool, error) {
	units, nextCursor, hasNext, err := s.repo.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var result []response.UnitResponse
	for _, unit := range units {
		result = append(result, *helper.MapUnit(&unit))
	}

	return result, nextCursor, hasNext, nil
}
