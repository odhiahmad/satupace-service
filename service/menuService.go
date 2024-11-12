package service

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/repository"
)

type MenuService interface {
	CreateMenu(menu request.MenuCreate)
	UpdateMenu(menu request.MenuUpdate)
	FindById(menuId int) response.MenuResponse
	FindAll() []response.MenuResponse
	Delete(menuId int)
}

type MenuRepository struct {
	MenuRepository repository.MenuRepository
	Validate       *validator.Validate
}

func NewMenuService(menuRepo repository.MenuRepository, validate *validator.Validate) MenuService {
	return &MenuRepository{
		MenuRepository: menuRepo,
		Validate:       validate,
	}
}

func (service *MenuRepository) CreateMenu(menu request.MenuCreate) {
	err := service.Validate.Struct(menu)
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	menuEntity := entity.Menu{
		Nama: menu.Nama,
	}

	service.MenuRepository.InsertMenu((menuEntity))
}

func (service *MenuRepository) UpdateMenu(menu request.MenuUpdate) {
	menuData, err := service.MenuRepository.FindById(menu.Id)
	helper.ErrorPanic(err)

	menuData.Nama = menu.Nama

	service.MenuRepository.UpdateMenu(menuData)
}

func (service *MenuRepository) FindById(menuId int) response.MenuResponse {
	menuData, err := service.MenuRepository.FindById(menuId)
	helper.ErrorPanic(err)

	tagResponse := response.MenuResponse{
		Id:   menuData.Id,
		Nama: menuData.Nama,
	}
	return tagResponse
}

// FindAll implements TagsService
func (t *MenuRepository) FindAll() []response.MenuResponse {
	result := t.MenuRepository.FindAll()

	var tags []response.MenuResponse
	for _, value := range result {
		tag := response.MenuResponse{
			Id:   value.Id,
			Nama: value.Nama,
		}
		tags = append(tags, tag)
	}

	return tags
}

func (t *MenuRepository) Delete(menuId int) {
	t.MenuRepository.Delete(menuId)
}
