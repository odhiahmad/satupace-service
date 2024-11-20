package service

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/mashingan/smapping"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/repository"
)

type UserService interface {
	CreateUser(user request.UserCreateDTO) entity.User
	UpdateUser(user request.UserUpdateDTO) entity.User
	IsDuplicateEmail(user string) bool
}

type userService struct {
	userRepository repository.UserRepository
	Validate       *validator.Validate
}

func NewUserService(userRepo repository.UserRepository, validate *validator.Validate) UserService {
	return &userService{
		userRepository: userRepo,
		Validate:       validate,
	}
}

func (service *userService) CreateUser(user request.UserCreateDTO) entity.User {
	userToCreate := entity.User{}
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}

	res := service.userRepository.InsertUser((userToCreate))
	return res
}

func (service *userService) UpdateUser(user request.UserUpdateDTO) entity.User {
	userToUpdate := entity.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	res := service.userRepository.UpdateUser((userToUpdate))
	return res
}

func (service *userService) IsDuplicateEmail(email string) bool {
	res := service.userRepository.IsDuplicateEmail(email)
	return !(res.Error == nil)
}
