package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type MembershipRepository interface {
	CreateMembership(membership entity.Membership) (entity.Membership, error)
	FindById(membershipId uuid.UUID) (membership entity.Membership, err error)
	FindActiveMembershipByBusinessID(businessId uuid.UUID) (*entity.Membership, error)
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

func (conn *membershipConnection) FindById(membershipId uuid.UUID) (memberships entity.Membership, err error) {
	var membership entity.Membership
	result := conn.db.Find(&membership, membershipId)
	if result != nil {
		return membership, nil
	} else {
		return membership, errors.New("tag is not found")
	}
}

func (conn *membershipConnection) FindActiveMembershipByBusinessID(businessId uuid.UUID) (*entity.Membership, error) {
	var membership entity.Membership
	err := conn.db.
		Where("business_id = ? AND is_active = ? AND end_date > ?", businessId, true, time.Now()).
		Order("end_date DESC").
		First(&membership).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("membership aktif tidak ditemukan")
		}
		return nil, err
	}

	return &membership, nil
}

func (conn *membershipConnection) FindAll() []entity.Membership {
	var membership []entity.Membership
	result := conn.db.Find(&membership)
	helper.ErrorPanic(result.Error)
	return membership
}
