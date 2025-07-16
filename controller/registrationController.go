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
}

type registrationController struct {
	registrationService service.RegistrationService
}

func NewRegistrationController(service service.RegistrationService) RegistrationController {
	return &registrationController{
		registrationService: service,
	}
}

func (c *registrationController) Register(ctx *gin.Context) {
	var req request.RegistrationRequest

	// Validasi binding JSON dari body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Gagal memproses data",
			"BAD_REQUEST",
			"request_body",
			err.Error(),
			helper.EmptyObj{},
		))
		return
	}

	// Panggil service untuk proses registrasi
	if err := c.registrationService.Register(req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Pendaftaran gagal",
			"REGISTER_FAILED",
			"register",
			err.Error(),
			helper.EmptyObj{},
		))
		return
	}

	// Response sukses tanpa data
	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Pendaftaran berhasil", nil))
}
