package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type AuthController interface {
	LoginBusiness(ctx *gin.Context)
	LoginPin(ctx *gin.Context)
	VerifyOTP(ctx *gin.Context)
	RetryOTP(ctx *gin.Context)
	RequestForgotPassword(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
	jwtService  service.JWTService
}

func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) LoginBusiness(ctx *gin.Context) {
	var loginDTO request.LoginUserBusinessDTO
	if err := ctx.ShouldBind(&loginDTO); err != nil {
		response := helper.BuildErrorResponse(
			"Input tidak valid",
			"VALIDATION_ERROR",
			"identifier",
			err.Error(),
			nil,
		)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	user, err := c.authService.VerifyCredentialBusiness(loginDTO.Identifier, loginDTO.Password)
	if err != nil {
		var response helper.ResponseError

		switch err {
		case helper.ErrInvalidPassword:
			response = helper.BuildErrorResponse(
				"Login gagal",
				"AUTH_INVALID_PASSWORD",
				"password",
				"Password yang Anda masukkan salah",
				nil,
			)
		case helper.ErrMembershipInactive:
			response = helper.BuildErrorResponse(
				"Login gagal",
				"AUTH_MEMBERSHIP_INACTIVE",
				"membership",
				"Membership Anda tidak aktif atau telah kedaluwarsa",
				nil,
			)
		case helper.ErrUserNotFound:
			response = helper.BuildErrorResponse(
				"Login gagal",
				"AUTH_USER_NOT_FOUND",
				"identifier",
				"Email atau nomor HP tidak ditemukan",
				nil,
			)
		default:
			response = helper.BuildErrorResponse(
				"Terjadi kesalahan",
				"AUTH_UNKNOWN_ERROR",
				"",
				err.Error(),
				nil,
			)
		}

		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	if !user.IsVerified {
		response := helper.BuildErrorResponse(
			"Akun belum diverifikasi",
			"AUTH_NOT_VERIFIED",
			"verification",
			"Silakan verifikasi akun Anda sebelum login",
			nil,
		)
		ctx.JSON(http.StatusForbidden, response)
		return
	}

	response := helper.BuildResponse(true, "Berhasil login", user)
	ctx.JSON(http.StatusOK, response)
}

func (c *authController) LoginPin(ctx *gin.Context) {
	var req request.PinLoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response := helper.BuildErrorResponse(
			"Input tidak valid",
			"VALIDATION_ERROR",
			"pin_login",
			err.Error(),
			nil,
		)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	authRes, err := c.authService.PinLogin(req)
	if err != nil {
		var response helper.ResponseError

		switch err {
		case helper.ErrUserNotFound:
			response = helper.BuildErrorResponse(
				"Login gagal",
				"AUTH_USER_NOT_FOUND",
				"user",
				"Nomor HP tidak ditemukan pada bisnis ini",
				nil,
			)
		case helper.ErrMembershipInactive:
			response = helper.BuildErrorResponse(
				"Login gagal",
				"AUTH_MEMBERSHIP_INACTIVE",
				"membership",
				"Membership Anda tidak aktif atau telah kedaluwarsa",
				nil,
			)
		default:
			response = helper.BuildErrorResponse(
				"Login gagal",
				"AUTH_PIN_ERROR",
				"pin",
				err.Error(),
				nil,
			)
		}

		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	response := helper.BuildResponse(true, "Berhasil login dengan PIN", authRes)
	ctx.JSON(http.StatusOK, response)
}

func (c *authController) VerifyOTP(ctx *gin.Context) {
	var req request.VerifyOTPRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid",
			"BAD_REQUEST",
			"request_body",
			err.Error(),
			helper.EmptyObj{},
		))
		return
	}

	userResponse, err := c.authService.VerifyOTPToken(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Verifikasi OTP gagal",
			"INVALID_TOKEN",
			"token",
			err.Error(),
			helper.EmptyObj{},
		))
		return
	}

	if req.IsResetPassword {
		ctx.JSON(http.StatusOK, helper.BuildResponse(true, "OTP berhasil diverifikasi", helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "OTP berhasil diverifikasi", userResponse))
}

func (c *authController) RetryOTP(ctx *gin.Context) {
	var req request.RetryOTPRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid",
			"bad_request",
			"request_body",
			err.Error(),
			helper.EmptyObj{},
		))
		return
	}

	if err := c.authService.RetryOTP(req); err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengirim ulang OTP",
			"retry_failed",
			"identifier",
			err.Error(),
			helper.EmptyObj{},
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "OTP berhasil dikirim ulang", helper.EmptyObj{}))
}

func (c *authController) RequestForgotPassword(ctx *gin.Context) {
	var req request.ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid",
			"bad_request",
			"request_body",
			err.Error(),
			helper.EmptyObj{},
		))
		return
	}

	if err := c.authService.RequestForgotPassword(req); err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengirim OTP reset password",
			"request_otp_failed",
			"identifier",
			err.Error(),
			helper.EmptyObj{},
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Kode OTP reset password berhasil dikirim", helper.EmptyObj{}))
}

func (c *authController) ResetPassword(ctx *gin.Context) {
	var req request.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid",
			"bad_request",
			"request_body",
			err.Error(),
			helper.EmptyObj{},
		))
		return
	}

	userResponse, err := c.authService.ResetPassword(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengatur ulang password",
			"reset_password_failed",
			"otp_or_identifier",
			err.Error(),
			helper.EmptyObj{},
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Password berhasil diatur ulang", userResponse))
}
