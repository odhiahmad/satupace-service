package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type ProductUnitController interface {
	CreateProductUnit(ctx *gin.Context)
	UpdateProductUnit(ctx *gin.Context)
	FindProductUnitById(ctx *gin.Context)
	FindProductUnitAll(ctx *gin.Context)
	DeleteProductUnit(ctx *gin.Context)
}

type productUnitController struct {
	productUnitService service.ProductUnitService
	jwtService         service.JWTService
}

func NewProductUnitController(productUnitService service.ProductUnitService, jwtService service.JWTService) ProductUnitController {
	return &productUnitController{
		productUnitService: productUnitService,
		jwtService:         jwtService,
	}
}

func (c *productUnitController) CreateProductUnit(ctx *gin.Context) {
	var productUnitCreate request.ProductUnitCreate
	err := ctx.ShouldBind(&productUnitCreate)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	c.productUnitService.CreateProductUnit(productUnitCreate)
	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusCreated, response)

}

func (c *productUnitController) UpdateProductUnit(ctx *gin.Context) {
	productUnitUpdate := request.ProductUnitUpdate{}
	err := ctx.ShouldBind(&productUnitUpdate)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	productUnitId := ctx.Param("productUnitId")

	id, err := strconv.Atoi(productUnitId)
	helper.ErrorPanic(err)
	productUnitUpdate.Id = id

	c.productUnitService.UpdateProductUnit(productUnitUpdate)

	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusCreated, response)

}

func (c *productUnitController) FindProductUnitAll(ctx *gin.Context) {
	productUnitResponse := c.productUnitService.FindAll()
	response := helper.BuildResponse(true, "!OK", productUnitResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *productUnitController) FindProductUnitById(ctx *gin.Context) {
	productUnitId := ctx.Param("productUnitId")
	id, err := strconv.Atoi(productUnitId)
	helper.ErrorPanic(err)

	productUnitResponse := c.productUnitService.FindById(id)

	response := helper.BuildResponse(true, "!OK", productUnitResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *productUnitController) DeleteProductUnit(ctx *gin.Context) {
	productUnitId := ctx.Param("productUnitId")
	id, err := strconv.Atoi(productUnitId)
	helper.ErrorPanic(err)

	c.productUnitService.Delete(id)

	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusOK, response)
}
