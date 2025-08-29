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

type TerminalService interface {
	Create(req request.TerminalRequest) (response.TerminalResponse, error)
	Update(id uuid.UUID, req request.TerminalRequest) (response.TerminalResponse, error)
	Delete(id uuid.UUID) error
	FindById(roleId uuid.UUID) response.TerminalResponse
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.TerminalResponse, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.TerminalResponse, string, bool, error)
}

type terminalService struct {
	repo     repository.TerminalRepository
	validate *validator.Validate
}

func NewTerminalService(repo repository.TerminalRepository, validate *validator.Validate) TerminalService {
	return &terminalService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *terminalService) Create(req request.TerminalRequest) (response.TerminalResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.TerminalResponse{}, err
	}

	terminal := entity.Terminal{
		BusinessId: req.BusinessId,
		Name:       strings.ToLower(req.Name),
		Location:   strings.ToLower(req.Location),
	}

	createdTerminal, err := s.repo.Create(terminal)
	if err != nil {
		return response.TerminalResponse{}, err
	}

	terminalResponse := mapper.MapTerminal(&createdTerminal)

	return *terminalResponse, nil
}

func (s *terminalService) Update(id uuid.UUID, req request.TerminalRequest) (response.TerminalResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return response.TerminalResponse{}, err
	}

	terminal := entity.Terminal{
		Id:       id,
		Name:     strings.ToLower(req.Name),
		Location: strings.ToLower(req.Location),
	}

	updatedTerminal, err := s.repo.Update(terminal)
	if err != nil {
		return response.TerminalResponse{}, err
	}

	terminalResponse := *mapper.MapTerminal(&updatedTerminal)

	return terminalResponse, nil
}

func (s *terminalService) Delete(id uuid.UUID) error {
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

func (s *terminalService) FindById(terminalId uuid.UUID) response.TerminalResponse {
	terminalData, err := s.repo.FindById(terminalId)
	helper.ErrorPanic(err)

	terminalResponse := mapper.MapTerminal(&terminalData)
	return *terminalResponse
}

func (s *terminalService) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]response.TerminalResponse, int64, error) {
	terminales, total, err := s.repo.FindWithPagination(businessId, pagination)
	if err != nil {
		return nil, 0, err
	}

	var result []response.TerminalResponse
	for _, terminal := range terminales {
		result = append(result, *mapper.MapTerminal(&terminal))
	}

	return result, total, nil
}

func (s *terminalService) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.TerminalResponse, string, bool, error) {
	terminals, nextCursor, hasNext, err := s.repo.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var result []response.TerminalResponse
	for _, terminal := range terminals {
		result = append(result, *mapper.MapTerminal(&terminal))
	}

	return result, nextCursor, hasNext, nil
}
