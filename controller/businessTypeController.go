package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type BusinessTypeController interface {
	CreateBusinessType(ctx *gin.Context)
	UpdateBusinessType(ctx *gin.Context)
	FindBusinessTypeById(ctx *gin.Context)
	FindBusinessTypeAll(ctx *gin.Context)
	DeleteBusinessType(ctx *gin.Context)
}

type businessTypeController struct {
	businessTypeService service.BusinessTypeService
	jwtService          service.JWTService
}

func NewBusinessTypeController(businessTypeService service.BusinessTypeService, jwtService service.JWTService) BusinessTypeController {
	return &businessTypeController{
		businessTypeService: businessTypeService,
		jwtService:          jwtService,
	}
}

func (c *businessTypeController) CreateBusinessType(ctx *gin.Context) {
	var businessTypeCreate request.BusinessTypeCreate
	err := ctx.ShouldBind(&businessTypeCreate)
	if err != nil {
		response := helper.BuildErrorResponse(
			"Failed to process request",
			"INVALID_REQUEST",
			"body",
			err.Error(),
			helper.EmptyObj{},
		)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	c.businessTypeService.CreateBusinessType(businessTypeCreate)
	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusCreated, response)
}

func (c *businessTypeController) UpdateBusinessType(ctx *gin.Context) {
	var businessTypeUpdate request.BusinessTypeUpdate
	err := ctx.ShouldBind(&businessTypeUpdate)
	if err != nil {
		response := helper.BuildErrorResponse(
			"Failed to process request",
			"INVALID_REQUEST",
			"body",
			err.Error(),
			helper.EmptyObj{},
		)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Konversi businessTypeId dari string ke int
	businessTypeIdParam := ctx.Param("businessTypeId")
	businessTypeId, err := strconv.Atoi(businessTypeIdParam)
	if err != nil {
		response := helper.BuildErrorResponse(
			"Invalid businessType ID",
			"INVALID_ID",
			"param",
			err.Error(),
			helper.EmptyObj{},
		)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Kirim ke service
	c.businessTypeService.UpdateBusinessType(businessTypeId, businessTypeUpdate)

	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusOK, response)
}

func (c *businessTypeController) FindBusinessTypeAll(ctx *gin.Context) {
	businessTypeResponse := c.businessTypeService.FindAll()
	response := helper.BuildResponse(true, "!OK", businessTypeResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *businessTypeController) FindBusinessTypeById(ctx *gin.Context) {
	businessTypeId := ctx.Param("businessTypeId")
	id, err := strconv.Atoi(businessTypeId)
	if err != nil {
		response := helper.BuildErrorResponse(
			"Invalid ID parameter",
			"INVALID_ID",
			"businessTypeId",
			err.Error(),
			helper.EmptyObj{},
		)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	businessTypeResponse := c.businessTypeService.FindById(id)
	response := helper.BuildResponse(true, "!OK", businessTypeResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *businessTypeController) DeleteBusinessType(ctx *gin.Context) {
	businessTypeId := ctx.Param("businessTypeId")
	id, err := strconv.Atoi(businessTypeId)
	if err != nil {
		response := helper.BuildErrorResponse(
			"Invalid ID parameter",
			"INVALID_ID",
			"businessTypeId",
			err.Error(),
			helper.EmptyObj{},
		)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	c.businessTypeService.Delete(id)

	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusOK, response)
}
