package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type RegistrationController interface {
	InsertRegistration(ctx *gin.Context)
}

type registrationController struct {
	registrationService service.RegistrationService
	jwtService          service.JWTService
}

func NewRegistrationController(registrationService service.RegistrationService, jwtService service.JWTService) RegistrationController {
	return &registrationController{
		registrationService: registrationService,
		jwtService:          jwtService,
	}
}

func (c *registrationController) InsertRegistration(ctx *gin.Context) {
	var registrationInsert request.Registration
	err := ctx.ShouldBind(&registrationInsert)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.registrationService.IsDuplicateEmail(registrationInsert.Email) {
		response := helper.BuildErrorResponse("Failed to process request", "Duplicate email", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
	} else {
		c.registrationService.Registration(registrationInsert)
		response := helper.BuildResponse(true, "!OK", nil)
		ctx.JSON(http.StatusCreated, response)
	}

}
