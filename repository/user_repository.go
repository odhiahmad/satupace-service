package repository

import (
	"log"

	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	InsertUser(user entity.User) entity.User
	UpdateUser(user entity.User) entity.User
	VerifyCredential(email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
	InsertRegistration(business entity.User)
}

type UserConnection struct {
	Db *gorm.DB
}

func NewUserRepository(Db *gorm.DB) UserRepository {
	return &UserConnection{Db: Db}
}

func (t *UserConnection) InsertUser(user entity.User) entity.User {
	user.Password = hashAndSalt([]byte(user.Password))
	t.Db.Save(&user)

	return user
}

func (t *UserConnection) UpdateUser(user entity.User) entity.User {

	if user.Password != "" {
		user.Password = hashAndSalt([]byte(user.Password))
	} else {
		var tempUser entity.User
		t.Db.Find(&tempUser, user.Email)
		user.Password = tempUser.Password
	}

	t.Db.Save(&user)

	return user
}

func (t *UserConnection) VerifyCredential(email string, password string) interface{} {
	var user entity.User
	res := t.Db.Where("email = ?", email).Take(&user)
	if res.Error == nil {
		return user
	}
	return nil
}

func (t *UserConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var user entity.User
	return t.Db.Where("email = ?", email).Take(&user)
}

func (t *UserConnection) InsertRegistration(user entity.User) {
	user.Password = hashAndSalt([]byte(user.Password))
	result := t.Db.Create(&user)

	helper.ErrorPanic(result.Error)
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}
	return string(hash)
}
