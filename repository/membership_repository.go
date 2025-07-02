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

type MembershipConnection struct {
	Db *gorm.DB
}

func NewMembershipRepository(Db *gorm.DB) MembershipRepository {
	return &MembershipConnection{Db: Db}
}

func (t *MembershipConnection) CreateMembership(membership entity.Membership) (entity.Membership, error) {
	result := t.Db.Create(&membership)
	helper.ErrorPanic(result.Error)

	return membership, result.Error
}

func (t *MembershipConnection) FindById(membershipId int) (memberships entity.Membership, err error) {
	var membership entity.Membership
	result := t.Db.Find(&membership, membershipId)
	if result != nil {
		return membership, nil
	} else {
		return membership, errors.New("tag is not found")
	}
}

func (t *MembershipConnection) FindAll() []entity.Membership {
	var membership []entity.Membership
	result := t.Db.Find(&membership)
	helper.ErrorPanic(result.Error)
	return membership
}
