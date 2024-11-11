package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type PerusahaanRepository interface {
	InsertPerusahaan(perusahaan entity.Perusahaan)
	UpdatePerusahaan(perusahaan entity.Perusahaan)
	FindById(perusahaanId int) (perusahaan entity.Perusahaan, err error)
	FindAll() []entity.Perusahaan
	Delete(perusahaanId int)
}

type PerusahaanConnection struct {
	Db *gorm.DB
}

func NewPerusahaanRepository(Db *gorm.DB) PerusahaanRepository {
	return &PerusahaanConnection{Db: Db}
}

func (t *PerusahaanConnection) InsertPerusahaan(perusahaan entity.Perusahaan) {
	result := t.Db.Create(&perusahaan)

	helper.ErrorPanic(result.Error)
}

func (t *PerusahaanConnection) UpdatePerusahaan(perusahaan entity.Perusahaan) {
	var updatePerusahaan = request.PerusahaanUpdateDTO{
		Id:     perusahaan.Id,
		Nama:   perusahaan.Nama,
		Alamat: perusahaan.Alamat,
		Lat:    perusahaan.Lat,
		Long:   perusahaan.Long,
		Logo:   perusahaan.Logo,
		Gambar: perusahaan.Gambar,
	}

	result := t.Db.Model(&perusahaan).Updates(updatePerusahaan)
	helper.ErrorPanic(result.Error)
}

func (t *PerusahaanConnection) FindById(perusahaanId int) (perusahaans entity.Perusahaan, err error) {
	var perusahaan entity.Perusahaan
	result := t.Db.Find(&perusahaan, perusahaanId)
	if result != nil {
		return perusahaan, nil
	} else {
		return perusahaan, errors.New("tag is not found")
	}
}

func (t *PerusahaanConnection) FindAll() []entity.Perusahaan {
	var perusahaan []entity.Perusahaan
	result := t.Db.Find(&perusahaan)
	helper.ErrorPanic(result.Error)
	return perusahaan
}

func (t *PerusahaanConnection) Delete(perusahaanId int) {
	var perusahaans entity.Perusahaan
	result := t.Db.Where("id = ?", perusahaanId).Delete(&perusahaans)
	helper.ErrorPanic(result.Error)
}
