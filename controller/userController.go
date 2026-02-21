package controller

import (
	"net/http"

	"run-sync/data/request"
	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindByEmail(ctx *gin.Context)
	FindByPhone(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	Delete(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
}

type userController struct {
	service service.UserService
}

func NewUserController(s service.UserService) UserController {
	return &userController{service: s}
}

func (c *userController) Create(ctx *gin.Context) {
	var req request.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Create(req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal membuat user", "CREATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "User berhasil dibuat", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *userController) Update(ctx *gin.Context) {
	userId, _ := uuid.Parse(ctx.Param("id"))
	var req request.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Update(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengubah user", "UPDATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "User berhasil diubah", result)
	ctx.JSON(http.StatusOK, response)
}

func (c *userController) FindById(ctx *gin.Context) {
	userId, _ := uuid.Parse(ctx.Param("id"))
	user, err := c.service.FindById(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data user", user)
	ctx.JSON(http.StatusOK, response)
}

func (c *userController) FindByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	user, err := c.service.FindByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data user", user)
	ctx.JSON(http.StatusOK, response)
}

func (c *userController) FindByPhone(ctx *gin.Context) {
	phone := ctx.Query("phone")
	user, err := c.service.FindByPhone(phone)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data user", user)
	ctx.JSON(http.StatusOK, response)
}

func (c *userController) FindAll(ctx *gin.Context) {
	users, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil semua data user", users)
	ctx.JSON(http.StatusOK, response)
}

func (c *userController) Delete(ctx *gin.Context) {
	userId, _ := uuid.Parse(ctx.Param("id"))
	err := c.service.Delete(userId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus user", "DELETE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "User berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
}

func (c *userController) ChangePassword(ctx *gin.Context) {
	userId, _ := uuid.Parse(ctx.Param("id"))
	var req request.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err := c.service.ChangePassword(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengubah password", "CHANGE_PASSWORD_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Password berhasil diubah", nil)
	ctx.JSON(http.StatusOK, response)
}
