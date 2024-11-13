package service

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/repository"
)

type UserService interface {
	CreateUser(user request.UserCreateDTO) entity.User
	UpdateUser(user request.UserUpdateDTO) entity.User
	IsDuplicateUsername(user string) bool
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
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

func (service *userService) IsDuplicateUsername(username string) bool {
	res := service.userRepository.IsDuplicateUsername(username)
	return !(res.Error == nil)
}
