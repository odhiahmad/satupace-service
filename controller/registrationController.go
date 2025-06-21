package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/service"
)

type RegistrationController interface {
	Register(ctx *gin.Context)
	CheckDuplicateEmail(ctx *gin.Context)
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

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	err := c.registrationService.Register(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
	})
}

func (c *registrationController) CheckDuplicateEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "email query parameter is required",
		})
		return
	}

	duplicate, err := c.registrationService.IsDuplicateEmail(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"email":         email,
		"is_duplicated": duplicate,
	})
}
