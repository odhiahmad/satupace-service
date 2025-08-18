package service

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/helper/mapper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type TableService interface {
	Create(req request.TableRequest) (response.TableResponse, error)
	Update(id uuid.UUID, req request.TableRequest) (response.TableResponse, error)
	Delete(id uuid.UUID) error
	FindById(roleId uuid.UUID) response.TableResponse
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.TableResponse, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.TableResponse, string, bool, error)
	GetActiveTables(businessId uuid.UUID) ([]response.TableWithTransactionsResponse, error)
}

type tableService struct {
	repo     repository.TableRepository
	validate *validator.Validate
}

func NewTableService(repo repository.TableRepository, validate *validator.Validate) TableService {
	return &tableService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *tableService) Create(req request.TableRequest) (response.TableResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.TableResponse{}, err
	}

	table := entity.Table{
		BusinessId: req.BusinessId,
		Number:     strings.ToLower(req.Number),
		Status:     req.Status,
	}

	createdTable, err := s.repo.Create(table)
	if err != nil {
		return response.TableResponse{}, err
	}

	tableResponse := mapper.MapTable(&createdTable)

	return *tableResponse, nil
}

func (s *tableService) Update(id uuid.UUID, req request.TableRequest) (response.TableResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.TableResponse{}, err
	}

	table := entity.Table{
		Id:         id,
		BusinessId: req.BusinessId,
		Number:     strings.ToLower(req.Number),
		Status:     req.Status,
	}

	updatedTable, err := s.repo.Update(table)
	if err != nil {
		return response.TableResponse{}, err
	}

	tableResponse := mapper.MapTable(&updatedTable)

	return *tableResponse, nil
}

func (s *tableService) Delete(id uuid.UUID) error {
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

func (s *tableService) FindById(tableId uuid.UUID) response.TableResponse {
	tableData, err := s.repo.FindById(tableId)
	helper.ErrorPanic(err)

	tableResponse := mapper.MapTable(&tableData)
	return *tableResponse
}

func (s *tableService) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.TableResponse, int64, error) {
	tablees, total, err := s.repo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.TableResponse
	for _, table := range tablees {
		result = append(result, *mapper.MapTable(&table))
	}

	return result, total, nil
}

func (s *tableService) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.TableResponse, string, bool, error) {
	tables, nextCursor, hasNext, err := s.repo.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var result []response.TableResponse
	for _, table := range tables {
		result = append(result, *mapper.MapTable(&table))
	}

	return result, nextCursor, hasNext, nil
}

func (s *tableService) GetActiveTables(businessId uuid.UUID) ([]response.TableWithTransactionsResponse, error) {
	tables, err := s.repo.GetActiveTables(businessId)
	if err != nil {
		return nil, err
	}

	var result []response.TableWithTransactionsResponse
	for _, tbl := range tables {
		result = append(result, *mapper.MapTableWithTransactions(&tbl))
	}

	return result, nil
}
