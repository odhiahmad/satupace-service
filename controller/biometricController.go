package controller

import (
	"net/http"

	"run-sync/data/request"
	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BiometricController interface {
	RegisterStart(ctx *gin.Context)
	RegisterFinish(ctx *gin.Context)
	LoginStart(ctx *gin.Context)
	LoginFinish(ctx *gin.Context)
	GetCredentials(ctx *gin.Context)
	DeleteCredential(ctx *gin.Context)
}

type biometricController struct {
	service service.BiometricService
}

func NewBiometricController(s service.BiometricService) BiometricController {
	return &biometricController{service: s}
}

func (c *biometricController) RegisterStart(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	var req request.BiometricRegisterStartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.RegisterStart(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal memulai registrasi biometrik", "BIOMETRIC_REGISTER_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Challenge biometrik berhasil dibuat", result)
	ctx.JSON(http.StatusOK, response)
}

func (c *biometricController) RegisterFinish(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	var req request.BiometricRegisterFinishRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.RegisterFinish(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menyelesaikan registrasi biometrik", "BIOMETRIC_REGISTER_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Biometrik berhasil didaftarkan", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *biometricController) LoginStart(ctx *gin.Context) {
	var req request.BiometricLoginStartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.LoginStart(req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal memulai login biometrik", "BIOMETRIC_LOGIN_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Challenge login biometrik berhasil dibuat", result)
	ctx.JSON(http.StatusOK, response)
}

func (c *biometricController) LoginFinish(ctx *gin.Context) {
	var req request.BiometricLoginFinishRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	user, token, err := c.service.LoginFinish(req)
	if err != nil {
		res := helper.BuildErrorResponse("Login biometrik gagal", "BIOMETRIC_LOGIN_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	response := helper.BuildResponse(true, "Login biometrik berhasil", map[string]interface{}{
		"user":  user,
		"token": token,
	})
	ctx.JSON(http.StatusOK, response)
}

func (c *biometricController) GetCredentials(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	credentials, err := c.service.GetCredentials(userId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil data biometrik", "BIOMETRIC_FETCH_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data biometrik", credentials)
	ctx.JSON(http.StatusOK, response)
}

func (c *biometricController) DeleteCredential(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	credentialId, _ := uuid.Parse(ctx.Param("id"))

	if err := c.service.DeleteCredential(userId, credentialId); err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus biometrik", "BIOMETRIC_DELETE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Biometrik berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
}
