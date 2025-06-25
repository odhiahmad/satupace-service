package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type UserBusinessRepository interface {
	InsertUserBusiness(userBusiness entity.UserBusiness) (entity.UserBusiness, error)
	FindById(userBusinessId int) (userBusiness entity.UserBusiness, err error)
	FindAll() []entity.UserBusiness
	Delete(userBusinessId int)
	IsDuplicateEmail(email string) bool
	VerifyCredentialBusiness(email string, password string) interface{}
	FindByEmailOrPhone(identifier string) (entity.UserBusiness, error)
}

type UserBusinessConnection struct {
	Db *gorm.DB
}

func NewUserBusinessRepository(Db *gorm.DB) UserBusinessRepository {
	return &UserBusinessConnection{Db: Db}
}

func (t *UserBusinessConnection) InsertUserBusiness(user entity.UserBusiness) (entity.UserBusiness, error) {
	result := t.Db.Create(&user)
	helper.ErrorPanic(result.Error)

	return user, result.Error
}

func (t *UserBusinessConnection) FindById(userBusinessId int) (userBusinesss entity.UserBusiness, err error) {
	var userBusiness entity.UserBusiness
	result := t.Db.Find(&userBusiness, userBusinessId)
	if result != nil {
		return userBusiness, nil
	} else {
		return userBusiness, errors.New("tag is not found")
	}
}

func (t *UserBusinessConnection) FindAll() []entity.UserBusiness {
	var userBusiness []entity.UserBusiness
	result := t.Db.Find(&userBusiness)
	helper.ErrorPanic(result.Error)
	return userBusiness
}

func (t *UserBusinessConnection) Delete(userBusinessId int) {
	var userBusinesss entity.UserBusiness
	result := t.Db.Where("id = ?", userBusinessId).Delete(&userBusinesss)
	helper.ErrorPanic(result.Error)
}

func (t *UserBusinessConnection) IsDuplicateEmail(email string) bool {
	var userBusiness entity.UserBusiness
	err := t.Db.Where("email = ?", email).Take(&userBusiness).Error

	// Perbaikan logika
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false // email belum ada → TIDAK duplikat
	}

	// Kalau tidak error atau error lain (misalnya koneksi), anggap duplikat
	return err == nil
}

func (t *UserBusinessConnection) VerifyCredentialBusiness(email string, password string) interface{} {
	var user entity.UserBusiness
	res := t.Db.Where("email = ?", email).Preload("Role").Preload("Business.BusinessType").Take(&user)
	if res.Error == nil {
		return user
	}
	return nil
}

func (t *UserBusinessConnection) FindByEmailOrPhone(identifier string) (entity.UserBusiness, error) {
	var user entity.UserBusiness
	res := t.Db.Where("email = ? OR phone_number = ?", identifier, identifier).
		Preload("Role").Preload("Business.BusinessType").
		First(&user)

	if res.Error != nil {
		return user, res.Error
	}
	return user, nil
}
