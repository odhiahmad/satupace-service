package service

import (
	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/repository"
)

type ProductPromoService interface {
	CreateMany(productId int, productVariantId int, businessId int, promos []request.ProductPromoCreate) error
	DeleteByProductId(productId int) error
}

type productPromoService struct {
	repo     repository.ProductPromoRepository
	validate *validator.Validate
}

func NewProductPromoService(repo repository.ProductPromoRepository) ProductPromoService {
	return &productPromoService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *productPromoService) CreateMany(productId int, productVariantId int, businessId int, promos []request.ProductPromoCreate) error {
	var productPromos []entity.ProductPromo

	for _, p := range promos {
		if err := s.validate.Struct(p); err != nil {
			return err
		}
		productPromos = append(productPromos, entity.ProductPromo{
			ProductId:        &productId,
			ProductVariantId: &productVariantId,
			PromoId:          p.PromoId,
			BusinessId:       businessId,
			MinQuantity:      p.MinQuantity,
		})
	}

	return s.repo.CreateMany(productPromos)
}

func (s *productPromoService) DeleteByProductId(productId int) error {
	return s.repo.DeleteByProductId(productId)
}
