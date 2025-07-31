package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type UserBusinessRepository interface {
	InsertUserBusiness(userBusiness entity.UserBusiness) (entity.UserBusiness, error)
	FindById(id uuid.UUID) (userBusiness entity.UserBusiness, err error)
	FindAll() []entity.UserBusiness
	Delete(userBusinessId uuid.UUID)
	IsDuplicateEmail(email string) bool
	VerifyCredentialBusiness(email string, password string) interface{}
	FindByEmailOrPhone(identifier string) (entity.UserBusiness, error)
	FindByVerificationToken(token string) (entity.UserBusiness, error)
	Update(user *entity.UserBusiness) error
}

type userBusinessConnection struct {
	db *gorm.DB
}

func NewUserBusinessRepository(db *gorm.DB) UserBusinessRepository {
	return &userBusinessConnection{db: db}
}

func (conn *userBusinessConnection) InsertUserBusiness(user entity.UserBusiness) (entity.UserBusiness, error) {
	result := conn.db.Create(&user)
	helper.ErrorPanic(result.Error)

	return user, result.Error
}

func (conn *userBusinessConnection) FindById(id uuid.UUID) (userBusinesss entity.UserBusiness, err error) {
	var userBusiness entity.UserBusiness
	result := conn.db.
		Preload("Role").
		Preload("Business").
		Preload("Business.BusinessType").
		Preload("Membership").
		Find(&userBusiness, id)
	if result != nil {
		return userBusiness, nil
	} else {
		return userBusiness, errors.New("tag is not found")
	}
}

func (conn *userBusinessConnection) FindAll() []entity.UserBusiness {
	var userBusiness []entity.UserBusiness
	result := conn.db.Find(&userBusiness)
	helper.ErrorPanic(result.Error)
	return userBusiness
}

func (conn *userBusinessConnection) Delete(userBusinessId uuid.UUID) {
	var userBusinesss entity.UserBusiness
	result := conn.db.Where("id = ?", userBusinessId).Delete(&userBusinesss)
	helper.ErrorPanic(result.Error)
}

func (conn *userBusinessConnection) IsDuplicateEmail(email string) bool {
	var userBusiness entity.UserBusiness
	err := conn.db.Where("email = ?", email).Take(&userBusiness).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	return err == nil
}

func (conn *userBusinessConnection) VerifyCredentialBusiness(email string, password string) interface{} {
	var user entity.UserBusiness
	res := conn.db.Where("email = ?", email).Preload("Role").Preload("Business.BusinessType").Take(&user)
	if res.Error == nil {
		return user
	}
	return nil
}

func (conn *userBusinessConnection) FindByEmailOrPhone(identifier string) (entity.UserBusiness, error) {
	var user entity.UserBusiness

	err := conn.db.
		Preload("Role").
		Preload("Business").
		Preload("Business.BusinessType").
		Preload("Membership").
		Where("email = ? OR phone_number = ? or pending_email = ?", identifier, identifier, identifier).
		First(&user).Error

	if err != nil {
		return entity.UserBusiness{}, err
	}

	return user, nil
}

func (conn *userBusinessConnection) FindByVerificationToken(token string) (entity.UserBusiness, error) {
	var user entity.UserBusiness
	err := conn.db.Where("verification_token = ?", token).First(&user).Error
	return user, err
}

func (r *userBusinessConnection) Update(user *entity.UserBusiness) error {
	if err := r.db.Save(user).Error; err != nil {
		return err
	}

	if user.Business.Id != uuid.Nil {
		if err := r.db.Save(user.Business).Error; err != nil {
			return err
		}
	}

	return nil
}
