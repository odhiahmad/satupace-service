package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type UserController interface {
	CreateUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

func NewUserController(userService service.UserService, jwtService service.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (c *userController) CreateUser(ctx *gin.Context) {
	var userCreateDTO request.UserCreateDTO
	errDTO := ctx.ShouldBind(&userCreateDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.userService.IsDuplicateUsername(userCreateDTO.Email) {
		response := helper.BuildErrorResponse("Failed to process request", "Duplicate username", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
	} else {
		createdUser := c.userService.CreateUser(userCreateDTO)
		response := helper.BuildResponse(true, "!OK", createdUser)
		ctx.JSON(http.StatusCreated, response)
	}
}

func (c *userController) UpdateUser(ctx *gin.Context) {
	var userUpdateDTO request.UserUpdateDTO
	errDTO := ctx.ShouldBind(&userUpdateDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.userService.IsDuplicateUsername(userUpdateDTO.Username) {
		response := helper.BuildErrorResponse("Failed to process request", "Duplicate username", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
	} else {
		updatedUser := c.userService.UpdateUser(userUpdateDTO)
		response := helper.BuildResponse(true, "!OK", updatedUser)
		ctx.JSON(http.StatusCreated, response)
	}
}
