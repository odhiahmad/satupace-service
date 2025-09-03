package controller

import (
	"net/http"

	"loka-kasir/helper"
	"loka-kasir/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HomeController interface {
	GetHome(ctx *gin.Context)
}

type homeController struct {
	homeService service.HomeService
	jwtService  service.JWTService
}

func NewHomeController(homeService service.HomeService, jwtService service.JWTService) HomeController {
	return &homeController{homeService: homeService, jwtService: jwtService}
}

func (c *homeController) GetHome(ctx *gin.Context) {
	businessIdStr := ctx.MustGet("business_id").(string)
	businessId, err := uuid.Parse(businessIdStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid business_id UUID"})
		return
	}

	dashboard, err := c.homeService.GetHome(businessId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data dashboard",
			"internal_error",
			"dashboard",
			err.Error(),
			nil,
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil data dashboard", dashboard))
}
