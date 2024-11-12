package repository

import (
	"errors"

	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"gorm.io/gorm"
)

type MenuRepository interface {
	InsertMenu(menu entity.Menu)
	UpdateMenu(menu entity.Menu)
	FindById(menuId int) (menu entity.Menu, err error)
	FindAll() []entity.Menu
	Delete(menuId int)
}

type MenuConnection struct {
	Db *gorm.DB
}

func NewMenuRepository(Db *gorm.DB) MenuRepository {
	return &MenuConnection{Db: Db}
}

func (t *MenuConnection) InsertMenu(menu entity.Menu) {
	result := t.Db.Create(&menu)

	helper.ErrorPanic(result.Error)
}

func (t *MenuConnection) UpdateMenu(menu entity.Menu) {
	var updateMenu = request.MenuUpdate{
		Id:   menu.Id,
		Nama: menu.Nama,
	}

	result := t.Db.Model(&menu).Updates(updateMenu)
	helper.ErrorPanic(result.Error)
}

func (t *MenuConnection) FindById(menuId int) (menus entity.Menu, err error) {
	var menu entity.Menu
	result := t.Db.Find(&menu, menuId)
	if result != nil {
		return menu, nil
	} else {
		return menu, errors.New("tag is not found")
	}
}

func (t *MenuConnection) FindAll() []entity.Menu {
	var menu []entity.Menu
	result := t.Db.Find(&menu)
	helper.ErrorPanic(result.Error)
	return menu
}

func (t *MenuConnection) Delete(menuId int) {
	var menus entity.Menu
	result := t.Db.Where("id = ?", menuId).Delete(&menus)
	helper.ErrorPanic(result.Error)
}
