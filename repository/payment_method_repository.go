package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/entity"
	"gorm.io/gorm"
)

type PaymentMethodRepository interface {
	InsertPaymentMethod(paymentMethod entity.PaymentMethod) (entity.PaymentMethod, error)
	UpdatePaymentMethod(paymentMethod entity.PaymentMethod) (entity.PaymentMethod, error)
	FindById(id int) (entity.PaymentMethod, error)
	FindAll() ([]entity.PaymentMethod, error)
	Delete(id int) error
}

type paymentMethodRepo struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) PaymentMethodRepository {
	return &paymentMethodRepo{db: db}
}

func (r *paymentMethodRepo) InsertPaymentMethod(paymentMethod entity.PaymentMethod) (entity.PaymentMethod, error) {
	err := r.db.Create(&paymentMethod).Error
	return paymentMethod, err
}

func (r *paymentMethodRepo) UpdatePaymentMethod(paymentMethod entity.PaymentMethod) (entity.PaymentMethod, error) {
	err := r.db.Save(&paymentMethod).Error
	return paymentMethod, err
}

func (r *paymentMethodRepo) FindById(id int) (entity.PaymentMethod, error) {
	var paymentMethod entity.PaymentMethod
	err := r.db.First(&paymentMethod, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.PaymentMethod{}, errors.New("payment method not found")
		}
		return entity.PaymentMethod{}, err
	}
	return paymentMethod, nil
}

func (r *paymentMethodRepo) FindAll() ([]entity.PaymentMethod, error) {
	var paymentMethods []entity.PaymentMethod
	err := r.db.Find(&paymentMethods).Error
	return paymentMethods, err
}

func (r *paymentMethodRepo) Delete(id int) error {
	return r.db.Delete(&entity.PaymentMethod{}, id).Error
}
