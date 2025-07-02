package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type PromoService interface {
	Create(req request.PromoCreate) (entity.Promo, error)
	Update(id int, req request.PromoUpdate) (entity.Promo, error)
	Delete(id int) error
	FindById(id int) (entity.Promo, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Promo, int64, error)
	SetIsActive(id int, active bool) error
}

type promoService struct {
	repo     repository.PromoRepository
	validate *validator.Validate
}

func NewPromoService(repo repository.PromoRepository, validate *validator.Validate) PromoService {
	return &promoService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *promoService) Create(req request.PromoCreate) (entity.Promo, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Promo{}, err
	}

	isPercentageVal := helper.DeterminePromoType(req.Amount)

	promo := entity.Promo{
		BusinessId:   req.BusinessId,
		Name:         req.Name,
		Description:  req.Description,
		Type:         req.Type,
		IsPercentage: isPercentageVal,
		Amount:       req.Amount,
		MinQuantity:  req.MinQuantity,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		IsActive:     req.IsActive,
	}

	// Simpan promo utama
	createdPromo, err := s.repo.Create(promo)
	if err != nil {
		return entity.Promo{}, err
	}

	// Simpan required products (many-to-many)
	if len(req.RequiredProductIds) > 0 {
		var requiredProducts []entity.Product
		for _, pid := range req.RequiredProductIds {
			requiredProducts = append(requiredProducts, entity.Product{Id: pid})
		}
		if err := s.repo.AppendRequiredProducts(&createdPromo, requiredProducts); err != nil {
			return entity.Promo{}, err
		}
	}

	return createdPromo, nil
}

func (s *promoService) Update(id int, req request.PromoUpdate) (entity.Promo, error) {
	// Validasi request
	if err := s.validate.Struct(req); err != nil {
		return entity.Promo{}, err
	}

	// Ambil data promo lama
	oldPromo, err := s.repo.FindById(id)
	if err != nil {
		return entity.Promo{}, err
	}

	// Tentukan tipe promo dari amount
	isPercentageVal := helper.DeterminePromoType(req.Amount)

	// Update field dasar
	oldPromo.Name = req.Name
	oldPromo.Description = req.Description
	oldPromo.Type = req.Type
	oldPromo.Amount = req.Amount
	oldPromo.IsPercentage = isPercentageVal
	oldPromo.MinQuantity = req.MinQuantity
	oldPromo.StartDate = req.StartDate
	oldPromo.EndDate = req.EndDate
	oldPromo.IsActive = req.IsActive

	// Simpan update promo ke database
	updatedPromo, err := s.repo.Update(oldPromo)
	if err != nil {
		return entity.Promo{}, err
	}

	// Mapping ulang RequiredProducts
	var requiredProducts []entity.Product
	for _, pid := range req.RequiredProductIds {
		requiredProducts = append(requiredProducts, entity.Product{Id: pid})
	}

	// Ganti isi tabel relasi many2many: promo_required_products
	if err := s.repo.ReplaceRequiredProducts(updatedPromo.Id, requiredProducts); err != nil {
		return entity.Promo{}, err
	}

	return updatedPromo, nil
}

func (s *promoService) Delete(id int) error {
	promo, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(promo)
}

func (s *promoService) FindById(id int) (entity.Promo, error) {
	return s.repo.FindById(id)
}

func (s *promoService) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Promo, int64, error) {
	return s.repo.FindWithPagination(businessId, pagination)
}

func (s *promoService) SetIsActive(id int, active bool) error {
	// Validasi keberadaan promo
	_, err := s.repo.FindById(id)
	if err != nil {
		return fmt.Errorf("promo tidak ditemukan: %w", err)
	}

	return s.repo.SetIsActive(id, active)
}
