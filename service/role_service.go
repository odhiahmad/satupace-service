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

type RoleService interface {
	CreateRole(role request.RoleCreate)
	UpdateRole(role request.RoleUpdate)
	FindById(roleId int) response.RoleResponse
	FindAll() []response.RoleResponse
	Delete(roleId int)
}

type RoleRepository struct {
	RoleRepository repository.RoleRepository
	Validate       *validator.Validate
}

func NewRoleService(roleRepo repository.RoleRepository, validate *validator.Validate) RoleService {
	return &RoleRepository{
		RoleRepository: roleRepo,
		Validate:       validate,
	}
}

func (service *RoleRepository) CreateRole(role request.RoleCreate) {
	err := service.Validate.Struct(role)
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	roleEntity := entity.Role{
		Name: role.Name,
	}

	service.RoleRepository.InsertRole((roleEntity))
}

func (service *RoleRepository) UpdateRole(role request.RoleUpdate) {
	roleData, err := service.RoleRepository.FindById(role.Id)
	helper.ErrorPanic(err)

	roleData.Name = role.Name

	service.RoleRepository.UpdateRole(roleData)
}

func (service *RoleRepository) FindById(roleId int) response.RoleResponse {
	roleData, err := service.RoleRepository.FindById(roleId)
	helper.ErrorPanic(err)

	tagResponse := response.RoleResponse{
		Id:   roleData.ID,
		Nama: roleData.Name,
	}
	return tagResponse
}

func (t *RoleRepository) FindAll() []response.RoleResponse {
	result := t.RoleRepository.FindAll()

	var tags []response.RoleResponse
	for _, value := range result {
		tag := response.RoleResponse{
			Id:   value.ID,
			Nama: value.Name,
		}
		tags = append(tags, tag)
	}

	return tags
}

func (t *RoleRepository) Delete(roleId int) {
	t.RoleRepository.Delete(roleId)
}
