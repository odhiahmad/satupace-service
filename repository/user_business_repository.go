package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type UserBusinessRepository interface {
	InsertUserBusiness(userBusiness entity.UserBusiness) (entity.UserBusiness, error)
	CreateEmployee(user *entity.UserBusiness) error
	FindById(id uuid.UUID) (userBusiness entity.UserBusiness, err error)
	FindAll() []entity.UserBusiness
	IsDuplicateEmail(email string) bool
	VerifyCredentialBusiness(email string, password string) interface{}
	FindByEmailOrPhone(identifier string) (entity.UserBusiness, error)
	FindByVerificationToken(token string) (entity.UserBusiness, error)
	Update(user *entity.UserBusiness) error
	FindByPhoneAndBusinessId(businessId uuid.UUID, phone string) (entity.UserBusiness, error)
	HasRelation(id uuid.UUID) (bool, error)
	SoftDelete(id uuid.UUID) error
	HardDelete(id uuid.UUID) error
	FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.UserBusiness, int64, error)
	FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.UserBusiness, string, bool, error)
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

func (r *userBusinessConnection) CreateEmployee(user *entity.UserBusiness) error {
	return r.db.Create(user).Error
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
	return r.db.Model(&entity.UserBusiness{}).
		Where("id = ?", user.Id).
		Updates(user).Error
}

func (conn *userBusinessConnection) FindByPhoneAndBusinessId(businessId uuid.UUID, phone string) (entity.UserBusiness, error) {
	var user entity.UserBusiness
	err := conn.db.Where("business_id = ? AND phone_number = ?", businessId, phone).First(&user).Error
	return user, err
}

func (conn *userBusinessConnection) HasRelation(id uuid.UUID) (bool, error) {
	var count int64

	if err := conn.db.Model(&entity.Shift{}).
		Where("cashier_id = ?", id).
		Count(&count).Error; err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	if err := conn.db.Model(&entity.Transaction{}).
		Where("cashier_id = ?", id).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (conn *userBusinessConnection) SoftDelete(id uuid.UUID) error {
	return conn.db.Delete(&entity.UserBusiness{}, id).Error
}

func (conn *userBusinessConnection) HardDelete(id uuid.UUID) error {
	return conn.db.Unscoped().Delete(&entity.UserBusiness{}, id).Error
}

func (conn *userBusinessConnection) FindWithPagination(businessId uuid.UUID, pagination request.Pagination) ([]entity.UserBusiness, int64, error) {
	var userBusiness []entity.UserBusiness
	var total int64

	baseQuery := conn.db.Model(&entity.UserBusiness{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		baseQuery = baseQuery.Where("name ILIKE ?", search)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	p := helper.Paginate(pagination, []string{"id", "name", "created_at", "updated_at"})

	_, _, err := p.Paginate(baseQuery, &userBusiness)
	if err != nil {
		return nil, 0, err
	}

	return userBusiness, total, nil
}

func (conn *userBusinessConnection) FindWithPaginationCursor(businessId uuid.UUID, pagination request.Pagination) ([]entity.UserBusiness, string, bool, error) {
	var userBusinesss []entity.UserBusiness

	query := conn.db.Model(&entity.UserBusiness{}).
		Where("business_id = ?", businessId)

	if pagination.Search != "" {
		search := "%" + pagination.Search + "%"
		query = query.Where("name ILIKE ?", search)
	}

	sortBy := pagination.SortBy
	if sortBy == "" {
		sortBy = "updated_at"
	}

	order := "ASC"
	if pagination.OrderBy == "desc" {
		order = "DESC"
	}

	if pagination.Cursor != "" {
		cursorID, err := helper.DecodeCursorID(pagination.Cursor)
		if err != nil {
			return nil, "", false, err
		}

		if order == "ASC" {
			query = query.Where("id > ?", cursorID)
		} else {
			query = query.Where("id < ?", cursorID)
		}
	}

	limit := pagination.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query = query.Order(fmt.Sprintf("%s %s", sortBy, order)).Limit(limit + 1)

	if err := query.Find(&userBusinesss).Error; err != nil {
		return nil, "", false, err
	}

	var nextCursor string
	hasNext := false

	if len(userBusinesss) > limit {
		last := userBusinesss[limit-1]
		nextCursor = helper.EncodeCursorID(last.Id.String())
		userBusinesss = userBusinesss[:limit]
		hasNext = true
	}

	return userBusinesss, nextCursor, hasNext, nil
}
