package repository

import (
	"errors"

	"loka-kasir/entity"

	"gorm.io/gorm"
)

type PaymentMethodRepository interface {
	InsertPaymentMethod(paymentMethod entity.PaymentMethod) (entity.PaymentMethod, error)
	UpdatePaymentMethod(paymentMethod entity.PaymentMethod) (entity.PaymentMethod, error)
	FindById(id int) (entity.PaymentMethod, error)
	FindAll() ([]entity.PaymentMethod, error)
	Delete(id int) error
}

type paymentMethodConnection struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) PaymentMethodRepository {
	return &paymentMethodConnection{db: db}
}

func (conn *paymentMethodConnection) InsertPaymentMethod(paymentMethod entity.PaymentMethod) (entity.PaymentMethod, error) {
	err := conn.db.Create(&paymentMethod).Error
	return paymentMethod, err
}

func (conn *paymentMethodConnection) UpdatePaymentMethod(paymentMethod entity.PaymentMethod) (entity.PaymentMethod, error) {
	err := conn.db.Save(&paymentMethod).Error
	return paymentMethod, err
}

func (conn *paymentMethodConnection) FindById(id int) (entity.PaymentMethod, error) {
	var paymentMethod entity.PaymentMethod
	err := conn.db.First(&paymentMethod, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.PaymentMethod{}, errors.New("payment method not found")
		}
		return entity.PaymentMethod{}, err
	}
	return paymentMethod, nil
}

func (conn *paymentMethodConnection) FindAll() ([]entity.PaymentMethod, error) {
	var paymentMethods []entity.PaymentMethod
	err := conn.db.Find(&paymentMethods).Error
	return paymentMethods, err
}

func (conn *paymentMethodConnection) Delete(id int) error {
	return conn.db.Delete(&entity.PaymentMethod{}, id).Error
}
