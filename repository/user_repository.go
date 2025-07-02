package repository

import (
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type UserRepository interface {
	InsertUser(user entity.User) entity.User
	UpdateUser(user entity.User) entity.User
	VerifyCredential(email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
	InsertRegistration(business entity.User)
}

type userConnection struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{db: db}
}

func (conn *userConnection) InsertUser(user entity.User) entity.User {
	user.Password = helper.HashAndSalt([]byte(user.Password))
	conn.db.Save(&user)

	return user
}

func (conn *userConnection) UpdateUser(user entity.User) entity.User {

	if user.Password != "" {
		user.Password = helper.HashAndSalt([]byte(user.Password))
	} else {
		var tempUser entity.User
		conn.db.Find(&tempUser, user.Email)
		user.Password = tempUser.Password
	}

	conn.db.Save(&user)

	return user
}

func (conn *userConnection) VerifyCredential(email string, password string) interface{} {
	var user entity.User
	res := conn.db.Where("email = ?", email).Preload("Business").Take(&user)
	if res.Error == nil {
		return user
	}
	return nil
}

func (conn *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var user entity.User
	return conn.db.Where("email = ?", email).Take(&user)
}

func (conn *userConnection) InsertRegistration(user entity.User) {
	user.Password = helper.HashAndSalt([]byte(user.Password))
	result := conn.db.Create(&user)

	helper.ErrorPanic(result.Error)
}
