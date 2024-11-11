package service

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type PerusahaanService interface {
	CreatePerusahaan(perusahaan request.PerusahaanCreateDTO)
	UpdatePerusahaan(perusahaan request.PerusahaanUpdateDTO)
	FindById(perusahaanId int) response.PerusahaanResponse
	FindAll() []response.PerusahaanResponse
	Delete(perusahaanId int)
}

type PerusahaanRepository struct {
	PerusahaanRepository repository.PerusahaanRepository
	Validate             *validator.Validate
}

func NewPerusahaanService(perusahaanRepo repository.PerusahaanRepository, validate *validator.Validate) PerusahaanService {
	return &PerusahaanRepository{
		PerusahaanRepository: perusahaanRepo,
		Validate:             validate,
	}
}

func (service *PerusahaanRepository) CreatePerusahaan(perusahaan request.PerusahaanCreateDTO) {
	err := service.Validate.Struct(perusahaan)
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	perusahaanEntity := entity.Perusahaan{
		Nama:     perusahaan.Nama,
		Alamat:   perusahaan.Alamat,
		Lat:      perusahaan.Lat,
		Long:     perusahaan.Long,
		Logo:     perusahaan.Logo,
		Gambar:   perusahaan.Gambar,
		IsActive: true,
	}

	service.PerusahaanRepository.InsertPerusahaan((perusahaanEntity))
}

func (service *PerusahaanRepository) UpdatePerusahaan(perusahaan request.PerusahaanUpdateDTO) {
	perusahaanData, err := service.PerusahaanRepository.FindById(perusahaan.Id)
	helper.ErrorPanic(err)

	perusahaanData.Nama = perusahaan.Nama
	perusahaanData.Alamat = perusahaan.Alamat
	perusahaanData.Lat = perusahaan.Lat
	perusahaanData.Long = perusahaan.Long
	perusahaanData.Logo = perusahaan.Logo
	perusahaanData.Gambar = perusahaan.Gambar

	service.PerusahaanRepository.UpdatePerusahaan(perusahaanData)
}

func (service *PerusahaanRepository) FindById(perusahaanId int) response.PerusahaanResponse {
	perusahaanData, err := service.PerusahaanRepository.FindById(perusahaanId)
	helper.ErrorPanic(err)

	tagResponse := response.PerusahaanResponse{
		Id:   perusahaanData.Id,
		Nama: perusahaanData.Nama,
	}
	return tagResponse
}

// FindAll implements TagsService
func (t *PerusahaanRepository) FindAll() []response.PerusahaanResponse {
	result := t.PerusahaanRepository.FindAll()

	var tags []response.PerusahaanResponse
	for _, value := range result {
		tag := response.PerusahaanResponse{
			Id:     value.Id,
			Nama:   value.Nama,
			Alamat: value.Alamat,
		}
		tags = append(tags, tag)
	}

	return tags
}

func (t *PerusahaanRepository) Delete(perusahaanId int) {
	t.PerusahaanRepository.Delete(perusahaanId)
}
