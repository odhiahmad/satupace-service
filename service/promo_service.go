package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/repository"
)

type PromoService interface {
	Create(req request.PromoCreate) (entity.Promo, error)
	Update(req request.PromoUpdate) (entity.Promo, error)
	Delete(id int) error
	FindById(id int) (entity.Promo, error)
	FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Promo, int64, error)
}

type promoService struct {
	repo             repository.PromoRepository
	productPromoRepo repository.ProductPromoRepository
	validate         *validator.Validate
}

func NewPromoService(repo repository.PromoRepository, productPromoRepo repository.ProductPromoRepository, validate *validator.Validate) PromoService {
	return &promoService{
		repo:             repo,
		productPromoRepo: productPromoRepo,
		validate:         validator.New(),
	}
}

func (s *promoService) Create(req request.PromoCreate) (entity.Promo, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Promo{}, err
	}

	promo := entity.Promo{
		BusinessId:  req.BusinessId,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Amount:      req.Amount,
		MinQuantity: req.MinQuantity,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		IsActive:    req.IsActive,
	}

	// Simpan promo utama
	createdPromo, err := s.repo.Create(promo)
	if err != nil {
		return entity.Promo{}, err
	}

	// Simpan relasi produk jika bukan global
	if !req.IsGlobal && len(req.ProductIds) > 0 {
		var productPromos []entity.ProductPromo
		for _, pid := range req.ProductIds {
			productPromos = append(productPromos, entity.ProductPromo{
				PromoId:     createdPromo.Id,
				ProductId:   pid,
				BusinessId:  req.BusinessId,
				MinQuantity: req.MinQuantity,
			})
		}
		err := s.productPromoRepo.CreateMany(productPromos)
		if err != nil {
			return entity.Promo{}, err
		}
	}

	return createdPromo, nil
}

func (s *promoService) Update(req request.PromoUpdate) (entity.Promo, error) {
	if err := s.validate.Struct(req); err != nil {
		return entity.Promo{}, err
	}

	// Ambil promo lama
	oldPromo, err := s.repo.FindById(req.Id)
	if err != nil {
		return entity.Promo{}, err
	}

	// Update field utama
	oldPromo.Name = req.Name
	oldPromo.Description = req.Description
	oldPromo.Type = req.Type
	oldPromo.Amount = req.Amount
	oldPromo.MinQuantity = req.MinQuantity
	oldPromo.StartDate = req.StartDate
	oldPromo.EndDate = req.EndDate
	oldPromo.IsActive = req.IsActive

	updatedPromo, err := s.repo.Update(oldPromo)
	if err != nil {
		return entity.Promo{}, err
	}

	// Hapus product promo lama
	_ = s.productPromoRepo.DeleteByPromoId(req.Id)

	// Simpan ulang product promos jika bukan global
	if !req.IsGlobal && len(req.ProductIds) > 0 {
		var productPromos []entity.ProductPromo
		for _, pid := range req.ProductIds {
			productPromos = append(productPromos, entity.ProductPromo{
				PromoId:     req.Id,
				ProductId:   pid,
				BusinessId:  oldPromo.BusinessId,
				MinQuantity: req.MinQuantity,
			})
		}
		err := s.productPromoRepo.CreateMany(productPromos)
		if err != nil {
			return entity.Promo{}, err
		}
	}

	return updatedPromo, nil
}

func (s *promoService) Delete(id int) error {
	promo, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	// Hapus relasi ProductPromo terlebih dahulu
	_ = s.productPromoRepo.DeleteByPromoId(id)

	return s.repo.Delete(promo)
}

func (s *promoService) FindById(id int) (entity.Promo, error) {
	return s.repo.FindById(id)
}

func (s *promoService) FindWithPagination(businessId int, pagination request.Pagination) ([]entity.Promo, int64, error) {
	return s.repo.FindWithPagination(businessId, pagination)
}
