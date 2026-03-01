package service

import (
	"errors"
	"log"
	"run-sync/data/request"
	"run-sync/data/response"
	"run-sync/entity"
	"run-sync/helper"
	"run-sync/repository"
	"time"

	"github.com/google/uuid"
)

func toUserPhotoResponse(photo *entity.UserPhoto) response.UserPhotoResponse {
	return response.UserPhotoResponse{
		Id:        photo.Id.String(),
		UserId:    photo.UserId.String(),
		Url:       photo.Url,
		Type:      photo.Type,
		IsPrimary: photo.IsPrimary,
		CreatedAt: photo.CreatedAt,
	}
}

type UserPhotoService interface {
	Create(userId uuid.UUID, req request.UploadUserPhotoRequest) (response.UserPhotoResponse, error)
	Update(id uuid.UUID, req request.UpdateUserPhotoRequest) (response.UserPhotoResponse, error)
	FindById(id uuid.UUID) (response.UserPhotoResponse, error)
	FindByUserId(userId uuid.UUID) ([]response.UserPhotoResponse, error)
	FindPrimaryPhoto(userId uuid.UUID) (response.UserPhotoResponse, error)
	Delete(id uuid.UUID) error
	VerifyFace(userId uuid.UUID, req request.FaceVerifyRequest) (response.FaceVerifyResponse, error)
}

type userPhotoService struct {
	repo     repository.UserPhotoRepository
	userRepo repository.UserRepository
}

func NewUserPhotoService(repo repository.UserPhotoRepository, userRepo repository.UserRepository) UserPhotoService {
	return &userPhotoService{repo: repo, userRepo: userRepo}
}

func (s *userPhotoService) Create(userId uuid.UUID, req request.UploadUserPhotoRequest) (response.UserPhotoResponse, error) {
	log.Printf("[UserPhoto] Create — type=%s isPrimary=%v imageLen=%d", req.Type, req.IsPrimary, len(req.Image))

	var faceWarning string
	// Jika tipe verification, coba validasi wajah — tapi jangan gagal jika tidak terdeteksi
	if req.Type == "verification" {
		if err := helper.DetectFrontFace(req.Image); err != nil {
			log.Printf("[UserPhoto] DetectFrontFace warning: %v", err)
			faceWarning = "Foto kurang jelas, pastikan wajah tampak jelas dan pencahayaan cukup. Foto tetap tersimpan namun verifikasi mungkin tidak akurat."
		}
	}

	// Upload image to Cloudinary
	imageUrl, err := helper.UploadBase64ToCloudinary(req.Image, "run-sync/photos")
	if err != nil {
		log.Printf("[UserPhoto] Cloudinary error: %v", err)
		return response.UserPhotoResponse{}, errors.New("gagal upload gambar ke Cloudinary: " + err.Error())
	}
	if imageUrl == "" {
		return response.UserPhotoResponse{}, errors.New("gambar tidak boleh kosong")
	}

	// If this is primary, set other photos to non-primary
	if req.IsPrimary {
		photos, _ := s.repo.FindByUserId(userId)
		for _, photo := range photos {
			photo.IsPrimary = false
			s.repo.Update(&photo)
		}
	}

	photo := entity.UserPhoto{
		Id:        uuid.New(),
		UserId:    userId,
		Url:       imageUrl,
		Type:      req.Type,
		IsPrimary: req.IsPrimary,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(&photo); err != nil {
		return response.UserPhotoResponse{}, err
	}

	result := toUserPhotoResponse(&photo)
	result.Warning = faceWarning
	return result, nil
}

func (s *userPhotoService) Update(id uuid.UUID, req request.UpdateUserPhotoRequest) (response.UserPhotoResponse, error) {
	photo, err := s.repo.FindById(id)
	if err != nil {
		return response.UserPhotoResponse{}, err
	}

	if req.Type != nil {
		photo.Type = *req.Type
	}
	if req.IsPrimary != nil {
		// If setting as primary, set others to non-primary
		if *req.IsPrimary {
			photos, _ := s.repo.FindByUserId(photo.UserId)
			for _, p := range photos {
				p.IsPrimary = false
				s.repo.Update(&p)
			}
		}
		photo.IsPrimary = *req.IsPrimary
	}

	if err := s.repo.Update(photo); err != nil {
		return response.UserPhotoResponse{}, err
	}

	return toUserPhotoResponse(photo), nil
}

func (s *userPhotoService) FindById(id uuid.UUID) (response.UserPhotoResponse, error) {
	photo, err := s.repo.FindById(id)
	if err != nil {
		return response.UserPhotoResponse{}, err
	}

	return toUserPhotoResponse(photo), nil
}

func (s *userPhotoService) FindByUserId(userId uuid.UUID) ([]response.UserPhotoResponse, error) {
	photos, err := s.repo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	var responses []response.UserPhotoResponse
	for _, photo := range photos {
		p := photo
		responses = append(responses, toUserPhotoResponse(&p))
	}

	return responses, nil
}

func (s *userPhotoService) FindPrimaryPhoto(userId uuid.UUID) (response.UserPhotoResponse, error) {
	photo, err := s.repo.FindPrimaryPhoto(userId)
	if err != nil {
		return response.UserPhotoResponse{}, err
	}

	return toUserPhotoResponse(photo), nil
}

func (s *userPhotoService) Delete(id uuid.UUID) error {
	photo, err := s.repo.FindById(id)
	if err != nil {
		return err
	}
	if photo.Type == "verification" {
		return errors.New("foto verifikasi tidak dapat dihapus")
	}
	return s.repo.Delete(id)
}

func (s *userPhotoService) VerifyFace(userId uuid.UUID, req request.FaceVerifyRequest) (response.FaceVerifyResponse, error) {
	// Ambil foto verifikasi milik user
	photo, err := s.repo.FindVerificationPhoto(userId)
	if err != nil {
		return response.FaceVerifyResponse{}, errors.New("foto verifikasi belum diunggah, silakan upload foto verifikasi terlebih dahulu")
	}

	// Bandingkan wajah kamera dengan foto verifikasi
	similarity, matched, err := helper.VerifyFaces(req.Image, photo.Url)
	if err != nil {
		return response.FaceVerifyResponse{}, errors.New("gagal melakukan verifikasi wajah: " + err.Error())
	}

	result := response.FaceVerifyResponse{
		Matched:    matched,
		Similarity: similarity,
	}
	if matched {
		p := toUserPhotoResponse(photo)
		result.Photo = &p

		// Set akun sebagai terverifikasi setelah wajah cocok
		user, err := s.userRepo.FindById(userId)
		if err != nil {
			return response.FaceVerifyResponse{}, errors.New("gagal mengambil data user: " + err.Error())
		}
		user.IsVerified = true
		if err := s.userRepo.Update(user); err != nil {
			return response.FaceVerifyResponse{}, errors.New("gagal memperbarui status verifikasi: " + err.Error())
		}
		result.IsVerified = true
	}
	return result, nil
}
