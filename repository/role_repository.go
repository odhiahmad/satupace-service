package repository

import (
	"errors"

	"loka-kasir/data/request"
	"loka-kasir/entity"
	"loka-kasir/helper"

	"gorm.io/gorm"
)

type RoleRepository interface {
	InsertRole(role entity.Role)
	UpdateRole(role entity.Role)
	FindById(roleId int) (role entity.Role, err error)
	FindAll() []entity.Role
	Delete(roleId int)
}

type roleConnection struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleConnection{db: db}
}

func (conn *roleConnection) InsertRole(role entity.Role) {
	result := conn.db.Create(&role)

	helper.ErrorPanic(result.Error)
}

func (conn *roleConnection) UpdateRole(role entity.Role) {
	var updateRole = request.RoleUpdate{
		Id:   role.Id,
		Name: role.Name,
	}

	result := conn.db.Model(&role).Updates(updateRole)
	helper.ErrorPanic(result.Error)
}

func (conn *roleConnection) FindById(roleId int) (roles entity.Role, err error) {
	var role entity.Role
	result := conn.db.Find(&role, roleId)
	if result != nil {
		return role, nil
	} else {
		return role, errors.New("tag is not found")
	}
}

func (conn *roleConnection) FindAll() []entity.Role {
	var role []entity.Role
	result := conn.db.Find(&role)
	helper.ErrorPanic(result.Error)
	return role
}

func (conn *roleConnection) Delete(roleId int) {
	var roles entity.Role
	result := conn.db.Where("id = ?", roleId).Delete(&roles)
	helper.ErrorPanic(result.Error)
}
