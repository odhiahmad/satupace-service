package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type RegistrationController interface {
	Register(ctx *gin.Context)
	CheckDuplicateEmail(ctx *gin.Context)
	VerifyOTP(ctx *gin.Context)
	RetryOTP(ctx *gin.Context)
}

type registrationController struct {
	registrationService service.RegistrationService
}

func NewRegistrationController(service service.RegistrationService) RegistrationController {
	return &registrationController{
		registrationService: service,
	}
}

// Register handles user registration
func (c *registrationController) Register(ctx *gin.Context) {
	var req request.RegistrationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse(
			"Gagal memproses data",
			"bad_request",
			"request_body",
			err.Error(),
			helper.EmptyObj{},
		)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	userResponse, err := c.registrationService.Register(req)
	if err != nil {
		statusCode := http.StatusBadRequest
		errorCode := "bad_request"
		field := "email"

		if err.Error() == "email sudah digunakan" {
			statusCode = http.StatusConflict // 409
			errorCode = "conflict"
		}

		res := helper.BuildErrorResponse(
			"Gagal melakukan registrasi",
			errorCode,
			field,
			err.Error(),
			helper.EmptyObj{},
		)
		ctx.JSON(statusCode, res)
		return
	}

	// ✅ Kirim data user saat registrasi berhasil
	res := helper.BuildResponse(true, "Registrasi berhasil", userResponse)
	ctx.JSON(http.StatusCreated, res)
}

// CheckDuplicateEmail handles checking if email is already registered
func (c *registrationController) CheckDuplicateEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		res := helper.BuildErrorResponse(
			"Parameter email wajib diisi",
			"bad_request",
			"email",
			"email query parameter is missing",
			helper.EmptyObj{},
		)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	duplicate, err := c.registrationService.IsDuplicateEmail(email)
	if err != nil {
		res := helper.BuildErrorResponse(
			"Gagal memeriksa email",
			"internal_error",
			"email",
			err.Error(),
			helper.EmptyObj{},
		)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	response := map[string]interface{}{
		"email":         email,
		"is_duplicated": duplicate,
	}

	res := helper.BuildResponse(true, "Pemeriksaan email berhasil", response)
	ctx.JSON(http.StatusOK, res)
}

// VerifyOTP memverifikasi kode OTP dari pengguna
func (c *registrationController) VerifyOTP(ctx *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Token string `json:"token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse(
			"Input tidak valid",
			"bad_request",
			"request_body",
			err.Error(),
			helper.EmptyObj{},
		)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.registrationService.VerifyOTPToken(req.Phone, req.Token); err != nil {
		res := helper.BuildErrorResponse(
			"Verifikasi OTP gagal",
			"invalid_token",
			"token",
			err.Error(),
			helper.EmptyObj{},
		)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "OTP berhasil diverifikasi", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

// RetryOTP mengirim ulang kode OTP ke pengguna
func (c *registrationController) RetryOTP(ctx *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse(
			"Input tidak valid",
			"bad_request",
			"request_body",
			err.Error(),
			helper.EmptyObj{},
		)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.registrationService.RetryOTP(req.Phone); err != nil {
		res := helper.BuildErrorResponse(
			"Gagal mengirim ulang OTP",
			"retry_failed",
			"phone",
			err.Error(),
			helper.EmptyObj{},
		)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := helper.BuildResponse(true, "OTP berhasil dikirim ulang", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}
