package service

import (
	"errors"
	"run-sync/data/request"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/helper"
	"run-sync/repository"
	"time"

	"github.com/google/uuid"
)

type UserService interface {
	Create(req request.CreateUserRequest) (response.UserDetailResponse, error)
	Update(id uuid.UUID, req request.UpdateUserRequest) (response.UserResponse, error)
	FindById(id uuid.UUID) (response.UserDetailResponse, error)
	FindByEmail(email string) (response.UserResponse, error)
	FindByPhone(phone string) (response.UserResponse, error)
	FindAll() ([]response.UserResponse, error)
	Delete(id uuid.UUID) error
	ChangePassword(id uuid.UUID, req request.ChangePasswordRequest) error
	Login(req request.LoginRequest) (response.UserResponse, error)
	VerifyAndActivate(phoneNumber string) (response.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(req request.CreateUserRequest) (response.UserDetailResponse, error) {
	// Check for duplicate email
	if req.Email != nil && s.repo.IsDuplicateEmail(*req.Email) {
		return response.UserDetailResponse{}, errors.New("email sudah terdaftar")
	}

	// Check for duplicate phone
	if s.repo.IsDuplicatePhone(req.PhoneNumber) {
		return response.UserDetailResponse{}, errors.New("nomor telepon sudah terdaftar")
	}

	hashedPassword := helper.HashPassword(req.Password)

	user := entity.User{
		Id:          uuid.New(),
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Gender:      req.Gender,
		Password:    hashedPassword,
		IsVerified:  false,
		IsActive:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(&user); err != nil {
		return response.UserDetailResponse{}, err
	}

	return response.UserDetailResponse{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		HasProfile:  user.HasProfile,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (s *userService) Update(id uuid.UUID, req request.UpdateUserRequest) (response.UserResponse, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		return response.UserResponse{}, err
	}

	if req.Name != nil {
		user.Name = req.Name
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}
	if req.PendingEmail != nil {
		user.PendingEmail = req.PendingEmail
	}
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(user); err != nil {
		return response.UserResponse{}, err
	}

	return response.UserResponse{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		HasProfile:  user.HasProfile,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (s *userService) FindById(id uuid.UUID) (response.UserDetailResponse, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		return response.UserDetailResponse{}, err
	}

	return response.UserDetailResponse{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		HasProfile:  user.HasProfile,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (s *userService) FindByEmail(email string) (response.UserResponse, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return response.UserResponse{}, err
	}

	return response.UserResponse{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		HasProfile:  user.HasProfile,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (s *userService) FindByPhone(phone string) (response.UserResponse, error) {
	user, err := s.repo.FindByPhone(phone)
	if err != nil {
		return response.UserResponse{}, err
	}

	return response.UserResponse{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		HasProfile:  user.HasProfile,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (s *userService) FindAll() ([]response.UserResponse, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []response.UserResponse
	for _, user := range users {
		responses = append(responses, response.UserResponse{
			Id:          user.Id.String(),
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Gender:      user.Gender,
			HasProfile:  user.HasProfile,
			IsVerified:  user.IsVerified,
			IsActive:    user.IsActive,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
	}

	return responses, nil
}

func (s *userService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *userService) ChangePassword(id uuid.UUID, req request.ChangePasswordRequest) error {
	user, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	if !helper.ComparePassword(user.Password, req.OldPassword) {
		return errors.New("password lama tidak sesuai")
	}

	user.Password = helper.HashPassword(req.NewPassword)
	user.UpdatedAt = time.Now()

	return s.repo.Update(user)
}

func (s *userService) Login(req request.LoginRequest) (response.UserResponse, error) {
	// Find user by email or phone number
	user, err := s.repo.FindByEmailOrPhone(req.Identifier)
	if err != nil {
		return response.UserResponse{}, errors.New("email/nomor telepon atau password salah")
	}

	// Check if account is suspended
	if user.IsSuspended {
		return response.UserResponse{}, errors.New("akun Anda telah disuspend")
	}

	// Check password
	if !helper.ComparePassword(user.Password, req.Password) {
		return response.UserResponse{}, errors.New("email/nomor telepon atau password salah")
	}

	return response.UserResponse{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		HasProfile:  user.HasProfile,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (s *userService) VerifyAndActivate(phoneNumber string) (response.UserResponse, error) {
	// Find user by phone
	user, err := s.repo.FindByPhone(phoneNumber)
	if err != nil {
		return response.UserResponse{}, errors.New("user tidak ditemukan")
	}

	// Update verification status
	user.IsVerified = true
	user.IsActive = true
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(user); err != nil {
		return response.UserResponse{}, err
	}

	return response.UserResponse{
		Id:          user.Id.String(),
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
		HasProfile:  user.HasProfile,
		IsVerified:  user.IsVerified,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}
