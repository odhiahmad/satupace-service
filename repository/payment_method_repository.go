package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type PaymentMethodRepository interface {
	InsertPaymentMethod(paymentMethod entity.PaymentMethod)
	UpdatePaymentMethod(paymentMethod entity.PaymentMethod)
	FindById(paymentMethodId int) (paymentMethod entity.PaymentMethod, err error)
	FindAll() []entity.PaymentMethod
	Delete(paymentMethodId int)
}

type PaymentMethodConnection struct {
	Db *gorm.DB
}

func NewPaymentMethodRepository(Db *gorm.DB) PaymentMethodRepository {
	return &PaymentMethodConnection{Db: Db}
}

func (t *PaymentMethodConnection) InsertPaymentMethod(paymentMethod entity.PaymentMethod) {
	result := t.Db.Create(&paymentMethod)

	helper.ErrorPanic(result.Error)
}

func (t *PaymentMethodConnection) UpdatePaymentMethod(paymentMethod entity.PaymentMethod) {
	var updatePaymentMethod = request.PaymentMethodUpdate{
		Id:   paymentMethod.Id,
		Name: paymentMethod.Name,
	}

	result := t.Db.Model(&paymentMethod).Updates(updatePaymentMethod)
	helper.ErrorPanic(result.Error)
}

func (t *PaymentMethodConnection) FindById(paymentMethodId int) (paymentMethods entity.PaymentMethod, err error) {
	var paymentMethod entity.PaymentMethod
	result := t.Db.Find(&paymentMethod, paymentMethodId)
	if result != nil {
		return paymentMethod, nil
	} else {
		return paymentMethod, errors.New("tag is not found")
	}
}

func (t *PaymentMethodConnection) FindAll() []entity.PaymentMethod {
	var paymentMethod []entity.PaymentMethod
	result := t.Db.Find(&paymentMethod)
	helper.ErrorPanic(result.Error)
	return paymentMethod
}

func (t *PaymentMethodConnection) Delete(paymentMethodId int) {
	var paymentMethods entity.PaymentMethod
	result := t.Db.Where("id = ?", paymentMethodId).Delete(&paymentMethods)
	helper.ErrorPanic(result.Error)
}
