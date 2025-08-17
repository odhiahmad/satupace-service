package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type UserBusinessController interface {
	FindById(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	ChangeEmail(ctx *gin.Context)
	ChangePhone(ctx *gin.Context)
}

type userBusinessController struct {
	service    service.UserBusinessService
	jwtService service.JWTService
}

func NewUserBusinessController(s service.UserBusinessService, jwtService service.JWTService) UserBusinessController {
	return &userBusinessController{service: s, jwtService: jwtService}
}

func (c *userBusinessController) FindById(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	user, err := c.service.FindById(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data userBusiness", user)
	ctx.JSON(http.StatusOK, response)
}

func (c *userBusinessController) ChangePassword(ctx *gin.Context) {
	id := ctx.MustGet("user_id").(uuid.UUID)
	var req request.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	req.Id = id
	if err := c.service.ChangePassword(req); err != nil {
		res := helper.BuildErrorResponse("Gagal mengubah password", "CHANGE_PASSWORD_FAILED", "password", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Password berhasil diubah", nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *userBusinessController) ChangeEmail(ctx *gin.Context) {
	id := ctx.MustGet("user_id").(uuid.UUID)
	var req request.ChangeEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	req.Id = id
	if err := c.service.ChangeEmail(req); err != nil {
		res := helper.BuildErrorResponse("Gagal mengubah email", "CHANGE_EMAIL_FAILED", "email", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Email berhasil diubah", nil)
	ctx.JSON(http.StatusOK, res)
}

func (c *userBusinessController) ChangePhone(ctx *gin.Context) {
	id := ctx.MustGet("user_id").(uuid.UUID)
	var req request.ChangePhoneRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	req.Id = id
	if err := c.service.ChangePhone(req); err != nil {
		res := helper.BuildErrorResponse("Gagal mengubah nomor HP", "CHANGE_PHONE_FAILED", "phone_number", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Nomor HP berhasil diubah, silakan verifikasi dengan OTP", nil)
	ctx.JSON(http.StatusOK, res)
}
