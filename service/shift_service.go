package service

import (
	"errors"
	"time"

	"loka-kasir/data/request"
	"loka-kasir/data/response"
	"loka-kasir/entity"
	"loka-kasir/helper/mapper"
	"loka-kasir/repository"

	"github.com/google/uuid"
)

type ShiftService interface {
	OpenShift(req request.OpenShiftRequest) error
	CloseShift(id uuid.UUID, req request.CloseShiftRequest) error
	GetActiveShift(terminalId string) (response.ShiftResponse, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.ShiftResponse, string, bool, error)
}

type shiftService struct {
	userRepo  repository.UserBusinessRepository
	shiftRepo repository.ShiftRepository
}

func NewShiftService(userRepo repository.UserBusinessRepository, shiftRepo repository.ShiftRepository) ShiftService {
	return &shiftService{userRepo, shiftRepo}
}

func (s *shiftService) OpenShift(req request.OpenShiftRequest) error {
	existingShift, _ := s.shiftRepo.FindOpenShiftByCashier(req.CashierId, req.TerminalId)
	if existingShift != nil {
		return errors.New("shift sudah dibuka")
	}

	shift := entity.Shift{
		Id:          uuid.New(),
		BusinessId:  req.BusinessId,
		TerminalId:  req.TerminalId,
		CashierId:   req.CashierId,
		OpenedAt:    time.Now(),
		OpeningCash: req.OpeningCash,
		Status:      "open",
		Notes:       req.Notes,
	}

	if err := s.shiftRepo.Create(&shift); err != nil {
		return err
	}

	return nil
}

func (s *shiftService) CloseShift(id uuid.UUID, req request.CloseShiftRequest) error {
	shift, err := s.shiftRepo.FindById(id)
	if err != nil {
		return errors.New("shift tidak ditemukan atau sudah ditutup")
	}

	now := time.Now()
	shift.ClosedAt = &now
	shift.ClosingCash = &req.ClosingCash
	shift.Status = "closed"
	shift.Notes = req.Notes

	if err := s.shiftRepo.Update(&shift); err != nil {
		return err
	}
	return nil
}

func (s *shiftService) GetActiveShift(terminalId string) (response.ShiftResponse, error) {
	id, err := uuid.Parse(terminalId)
	if err != nil {
		return response.ShiftResponse{}, errors.New("invalid terminal id")
	}

	shift, err := s.shiftRepo.GetActiveShiftByTerminal(id)
	if err != nil {
		return response.ShiftResponse{}, errors.New("no active shift found for this terminal")
	}

	shiftResponse := mapper.MapShift(shift)
	return *shiftResponse, nil
}

func (s *shiftService) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]response.ShiftResponse, string, bool, error) {
	shifts, nextCursor, hasNext, err := s.shiftRepo.FindWithPaginationCursor(businessId, pagination)
	if err != nil {
		return nil, "", false, err
	}

	var result []response.ShiftResponse
	for _, shift := range shifts {
		result = append(result, *mapper.MapShift(&shift))
	}

	return result, nextCursor, hasNext, nil
}
