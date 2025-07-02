package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type MembershipRepository interface {
	CreateMembership(membership entity.Membership) (entity.Membership, error)
	FindById(membershipId int) (membership entity.Membership, err error)
	FindAll() []entity.Membership
}

type membershipConnection struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) MembershipRepository {
	return &membershipConnection{db: db}
}

func (conn *membershipConnection) CreateMembership(membership entity.Membership) (entity.Membership, error) {
	result := conn.db.Create(&membership)
	helper.ErrorPanic(result.Error)

	return membership, result.Error
}

func (conn *membershipConnection) FindById(membershipId int) (memberships entity.Membership, err error) {
	var membership entity.Membership
	result := conn.db.Find(&membership, membershipId)
	if result != nil {
		return membership, nil
	} else {
		return membership, errors.New("tag is not found")
	}
}

func (conn *membershipConnection) FindAll() []entity.Membership {
	var membership []entity.Membership
	result := conn.db.Find(&membership)
	helper.ErrorPanic(result.Error)
	return membership
}
