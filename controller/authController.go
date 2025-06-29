package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/entity"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type AuthController interface {
	Login(ctx *gin.Context)
	LoginBusiness(ctx *gin.Context)
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

func (c *authController) Login(ctx *gin.Context) {
	var loginDTO request.LoginUserDTO
	if err := ctx.ShouldBind(&loginDTO); err != nil {
		response := helper.BuildErrorResponse(
			"Input tidak valid",
			"VALIDATION_ERROR",
			"email",
			err.Error(),
			nil,
		)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if user, ok := authResult.(entity.User); ok {
		token := c.jwtService.GenerateToken(user.Id)
		user.Token = token
		response := helper.BuildResponse(true, "Berhasil login", user)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := helper.BuildErrorResponse(
		"Login gagal",
		"AUTH_INVALID_CREDENTIALS",
		"email",
		"Email atau password tidak valid",
		nil,
	)
	ctx.JSON(http.StatusUnauthorized, response)
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
		response := helper.BuildErrorResponse(
			"Login gagal",
			"AUTH_INVALID_CREDENTIALS",
			"identifier",
			"Email/Nomor HP atau password tidak valid",
			nil,
		)
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	// Tidak perlu generate token lagi, sudah di dalam service
	response := helper.BuildResponse(true, "Berhasil login", user)
	ctx.JSON(http.StatusOK, response)
}
