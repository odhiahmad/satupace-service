package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type PerusahaanController interface {
	CreatePerusahaan(ctx *gin.Context)
	UpdatePerusahaan(ctx *gin.Context)
	FindPerusahaanById(ctx *gin.Context)
	FindPerusahaanAll(ctx *gin.Context)
	DeletePerusahaan(ctx *gin.Context)
}

type perusahaanController struct {
	perusahaanService service.PerusahaanService
	jwtService        service.JWTService
}

func NewPerusahaanController(perusahaanService service.PerusahaanService, jwtService service.JWTService) PerusahaanController {
	return &perusahaanController{
		perusahaanService: perusahaanService,
		jwtService:        jwtService,
	}
}

func (c *perusahaanController) CreatePerusahaan(ctx *gin.Context) {
	var perusahaanCreateDTO request.PerusahaanCreateDTO
	errDTO := ctx.ShouldBind(&perusahaanCreateDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	c.perusahaanService.CreatePerusahaan(perusahaanCreateDTO)
	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusCreated, response)

}

func (c *perusahaanController) UpdatePerusahaan(ctx *gin.Context) {
	perusahaanUpdateDTO := request.PerusahaanUpdateDTO{}
	errDTO := ctx.ShouldBind(&perusahaanUpdateDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	perusahaanId := ctx.Param("perusahaanId")

	id, err := strconv.Atoi(perusahaanId)
	helper.ErrorPanic(err)
	perusahaanUpdateDTO.Id = id

	c.perusahaanService.UpdatePerusahaan(perusahaanUpdateDTO)

	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusCreated, response)

}

func (c *perusahaanController) FindPerusahaanAll(ctx *gin.Context) {
	perusahaanResponse := c.perusahaanService.FindAll()
	response := helper.BuildResponse(true, "!OK", perusahaanResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *perusahaanController) FindPerusahaanById(ctx *gin.Context) {
	perusahaanId := ctx.Param("perusahaanId")
	id, err := strconv.Atoi(perusahaanId)
	helper.ErrorPanic(err)

	perusahaanResponse := c.perusahaanService.FindById(id)

	response := helper.BuildResponse(true, "!OK", perusahaanResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *perusahaanController) DeletePerusahaan(ctx *gin.Context) {
	perusahaanId := ctx.Param("perusahaanId")
	id, err := strconv.Atoi(perusahaanId)
	helper.ErrorPanic(err)

	c.perusahaanService.Delete(id)

	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusOK, response)
}
