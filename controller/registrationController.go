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

// Register handles user registration
func (c *registrationController) Register(ctx *gin.Context) {
	var req request.RegistrationRequest

	// Bind & validate JSON input
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid JSON: " + err.Error(),
		})
		return
	}

	// Call service
	if err := c.registrationService.Register(req); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "email sudah digunakan" {
			statusCode = http.StatusConflict // 409
		}
		ctx.JSON(statusCode, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Success response
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Registration successful",
	})
}

// CheckDuplicateEmail handles checking if email is already registered
func (c *registrationController) CheckDuplicateEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Query parameter 'email' is required",
		})
		return
	}

	duplicate, err := c.registrationService.IsDuplicateEmail(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"email":         email,
		"is_duplicated": duplicate,
	})
}
