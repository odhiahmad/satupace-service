package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/repository"
)

type ShiftService interface {
	OpenShift(req request.OpenShiftRequest) (*entity.Shift, error)
	CloseShift(req request.CloseShiftRequest) (*entity.Shift, error)
	GetActiveShift(terminalId string) (*entity.Shift, error)
}

type shiftService struct {
	userRepo  repository.UserBusinessRepository
	shiftRepo repository.ShiftRepository
}

func NewShiftService(userRepo repository.UserBusinessRepository, shiftRepo repository.ShiftRepository) ShiftService {
	return &shiftService{userRepo, shiftRepo}
}

func (s *shiftService) OpenShift(req request.OpenShiftRequest) (*entity.Shift, error) {
	existingShift, _ := s.shiftRepo.FindOpenShiftByCashier(req.CashierId, req.TerminalId)
	if existingShift != nil {
		return nil, errors.New("shift sudah dibuka")
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
		return nil, err
	}
	return &shift, nil
}

func (s *shiftService) CloseShift(req request.CloseShiftRequest) (*entity.Shift, error) {
	shift, err := s.shiftRepo.FindOpenShiftByCashier(req.CashierId, req.TerminalId)
	if err != nil {
		return nil, errors.New("shift tidak ditemukan atau sudah ditutup")
	}

	now := time.Now()
	shift.ClosedAt = &now
	shift.ClosingCash = &req.ClosingCash
	shift.Status = "closed"
	shift.Notes = req.Notes

	if err := s.shiftRepo.Update(shift); err != nil {
		return nil, err
	}
	return shift, nil
}

func (s *shiftService) GetActiveShift(terminalId string) (*entity.Shift, error) {
	id, err := uuid.Parse(terminalId)
	if err != nil {
		return nil, errors.New("invalid terminal id")
	}

	shift, err := s.shiftRepo.GetActiveShiftByTerminal(id)
	if err != nil {
		return nil, errors.New("no active shift found for this terminal")
	}

	return shift, nil
}
