package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type UserBusinessRepository interface {
	InsertUserBusiness(userBusiness entity.UserBusiness) entity.UserBusiness
	FindById(userBusinessId int) (userBusiness entity.UserBusiness, err error)
	FindAll() []entity.UserBusiness
	Delete(userBusinessId int)
	IsDuplicateEmail(email string) (tx *gorm.DB)
	VerifyCredentialBusiness(email string, password string) interface{}
}

type UserBusinessConnection struct {
	Db *gorm.DB
}

func NewUserBusinessRepository(Db *gorm.DB) UserBusinessRepository {
	return &UserBusinessConnection{Db: Db}
}

func (t *UserBusinessConnection) InsertUserBusiness(userBusiness entity.UserBusiness) entity.UserBusiness {
	result := t.Db.Create(&userBusiness)
	helper.ErrorPanic(result.Error)

	return userBusiness
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

func (t *UserBusinessConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var userBusiness entity.UserBusiness
	return t.Db.Where("email = ?", email).Take(&userBusiness)
}

func (t *UserBusinessConnection) VerifyCredentialBusiness(email string, password string) interface{} {
	var user entity.UserBusiness
	res := t.Db.Where("email = ?", email).Preload("Role").Preload("Business.BusinessType").Take(&user)
	if res.Error == nil {
		return user
	}
	return nil
}
