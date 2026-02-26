package service

import (
	"errors"
	"run-sync/data/request"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/repository"
	"time"

	"github.com/google/uuid"
)

var dayNames = [7]string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}

type RunGroupScheduleService interface {
	Create(groupId uuid.UUID, req request.CreateRunGroupScheduleRequest) (response.RunGroupScheduleResponse, error)
	Update(scheduleId uuid.UUID, req request.UpdateRunGroupScheduleRequest) (response.RunGroupScheduleResponse, error)
	FindByGroupId(groupId uuid.UUID) ([]response.RunGroupScheduleResponse, error)
	Delete(scheduleId uuid.UUID) error
}

type runGroupScheduleService struct {
	repo      repository.RunGroupScheduleRepository
	groupRepo repository.RunGroupRepository
}

func NewRunGroupScheduleService(repo repository.RunGroupScheduleRepository, groupRepo repository.RunGroupRepository) RunGroupScheduleService {
	return &runGroupScheduleService{repo: repo, groupRepo: groupRepo}
}

func (s *runGroupScheduleService) Create(groupId uuid.UUID, req request.CreateRunGroupScheduleRequest) (response.RunGroupScheduleResponse, error) {
	// Pastikan group ada
	if _, err := s.groupRepo.FindById(groupId); err != nil {
		return response.RunGroupScheduleResponse{}, errors.New("grup lari tidak ditemukan")
	}

	// Maksimal 3 jadwal aktif per group
	count, err := s.repo.CountByGroupId(groupId)
	if err != nil {
		return response.RunGroupScheduleResponse{}, err
	}
	if count >= 3 {
		return response.RunGroupScheduleResponse{}, errors.New("maksimal 3 jadwal per grup lari")
	}

	schedule := entity.RunGroupSchedule{
		Id:        uuid.New(),
		GroupId:   groupId,
		DayOfWeek: req.DayOfWeek,
		StartTime: req.StartTime,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(&schedule); err != nil {
		return response.RunGroupScheduleResponse{}, err
	}

	return mapScheduleEntityToResponse(&schedule), nil
}

func (s *runGroupScheduleService) Update(scheduleId uuid.UUID, req request.UpdateRunGroupScheduleRequest) (response.RunGroupScheduleResponse, error) {
	schedule, err := s.repo.FindById(scheduleId)
	if err != nil {
		return response.RunGroupScheduleResponse{}, errors.New("jadwal tidak ditemukan")
	}

	if req.DayOfWeek != nil {
		schedule.DayOfWeek = *req.DayOfWeek
	}
	if req.StartTime != nil {
		schedule.StartTime = *req.StartTime
	}
	if req.IsActive != nil {
		schedule.IsActive = *req.IsActive
	}
	schedule.UpdatedAt = time.Now()

	if err := s.repo.Update(schedule); err != nil {
		return response.RunGroupScheduleResponse{}, err
	}

	return mapScheduleEntityToResponse(schedule), nil
}

func (s *runGroupScheduleService) FindByGroupId(groupId uuid.UUID) ([]response.RunGroupScheduleResponse, error) {
	schedules, err := s.repo.FindByGroupId(groupId)
	if err != nil {
		return nil, err
	}

	var responses []response.RunGroupScheduleResponse
	for _, sch := range schedules {
		responses = append(responses, mapScheduleEntityToResponse(&sch))
	}

	return responses, nil
}

func (s *runGroupScheduleService) Delete(scheduleId uuid.UUID) error {
	if _, err := s.repo.FindById(scheduleId); err != nil {
		return errors.New("jadwal tidak ditemukan")
	}
	return s.repo.Delete(scheduleId)
}

func mapScheduleEntityToResponse(s *entity.RunGroupSchedule) response.RunGroupScheduleResponse {
	dayName := ""
	if s.DayOfWeek >= 0 && s.DayOfWeek <= 6 {
		dayName = dayNames[s.DayOfWeek]
	}
	return response.RunGroupScheduleResponse{
		Id:        s.Id.String(),
		GroupId:   s.GroupId.String(),
		DayOfWeek: s.DayOfWeek,
		DayName:   dayName,
		StartTime: s.StartTime,
		IsActive:  s.IsActive,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
