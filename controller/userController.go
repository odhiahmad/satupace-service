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
		response := helper.BuildErrorResponse(
			"Permintaan tidak valid",
			"BAD_REQUEST",
			"body",
			errDTO.Error(),
			nil,
		)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.userService.IsDuplicateEmail(userCreateDTO.Email) {
		response := helper.BuildErrorResponse(
			"Email sudah digunakan",
			"DUPLICATE_EMAIL",
			"email",
			"Alamat email telah terdaftar sebelumnya",
			nil,
		)
		ctx.JSON(http.StatusConflict, response)
		return
	}

	createdUser := c.userService.CreateUser(userCreateDTO)
	response := helper.BuildResponse(true, "Berhasil membuat pengguna", createdUser)
	ctx.JSON(http.StatusCreated, response)
}

func (c *userController) UpdateUser(ctx *gin.Context) {
	var userUpdateDTO request.UserUpdateDTO
	errDTO := ctx.ShouldBind(&userUpdateDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse(
			"Permintaan tidak valid",
			"BAD_REQUEST",
			"body",
			errDTO.Error(),
			nil,
		)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.userService.IsDuplicateEmail(userUpdateDTO.Email) {
		response := helper.BuildErrorResponse(
			"Email sudah digunakan",
			"DUPLICATE_EMAIL",
			"email",
			"Alamat email telah terdaftar sebelumnya",
			nil,
		)
		ctx.JSON(http.StatusConflict, response)
		return
	}

	updatedUser := c.userService.UpdateUser(userUpdateDTO)
	response := helper.BuildResponse(true, "Berhasil memperbarui pengguna", updatedUser)
	ctx.JSON(http.StatusOK, response)
}
