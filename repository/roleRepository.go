package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type RoleRepository interface {
	InsertRole(role entity.Role)
	UpdateRole(role entity.Role)
	FindById(roleId int) (role entity.Role, err error)
	FindAll() []entity.Role
	Delete(roleId int)
}

type RoleConnection struct {
	Db *gorm.DB
}

func NewRoleRepository(Db *gorm.DB) RoleRepository {
	return &RoleConnection{Db: Db}
}

func (t *RoleConnection) InsertRole(role entity.Role) {
	result := t.Db.Create(&role)

	helper.ErrorPanic(result.Error)
}

func (t *RoleConnection) UpdateRole(role entity.Role) {
	var updateRole = request.RoleUpdate{
		Id:   role.Id,
		Nama: role.Nama,
	}

	result := t.Db.Model(&role).Updates(updateRole)
	helper.ErrorPanic(result.Error)
}

func (t *RoleConnection) FindById(roleId int) (roles entity.Role, err error) {
	var role entity.Role
	result := t.Db.Find(&role, roleId)
	if result != nil {
		return role, nil
	} else {
		return role, errors.New("tag is not found")
	}
}

func (t *RoleConnection) FindAll() []entity.Role {
	var role []entity.Role
	result := t.Db.Find(&role)
	helper.ErrorPanic(result.Error)
	return role
}

func (t *RoleConnection) Delete(roleId int) {
	var roles entity.Role
	result := t.Db.Where("id = ?", roleId).Delete(&roles)
	helper.ErrorPanic(result.Error)
}
