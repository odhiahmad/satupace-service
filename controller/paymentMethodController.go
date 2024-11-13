package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type PaymentMethodController interface {
	CreatePaymentMethod(ctx *gin.Context)
	UpdatePaymentMethod(ctx *gin.Context)
	FindPaymentMethodById(ctx *gin.Context)
	FindPaymentMethodAll(ctx *gin.Context)
	DeletePaymentMethod(ctx *gin.Context)
}

type paymentMethodController struct {
	paymentMethodService service.PaymentMethodService
	jwtService           service.JWTService
}

func NewPaymentMethodController(paymentMethodService service.PaymentMethodService, jwtService service.JWTService) PaymentMethodController {
	return &paymentMethodController{
		paymentMethodService: paymentMethodService,
		jwtService:           jwtService,
	}
}

func (c *paymentMethodController) CreatePaymentMethod(ctx *gin.Context) {
	var paymentMethodCreate request.PaymentMethodCreate
	err := ctx.ShouldBind(&paymentMethodCreate)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	c.paymentMethodService.CreatePaymentMethod(paymentMethodCreate)
	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusCreated, response)

}

func (c *paymentMethodController) UpdatePaymentMethod(ctx *gin.Context) {
	paymentMethodUpdate := request.PaymentMethodUpdate{}
	err := ctx.ShouldBind(&paymentMethodUpdate)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	paymentMethodId := ctx.Param("paymentMethodId")

	id, err := strconv.Atoi(paymentMethodId)
	helper.ErrorPanic(err)
	paymentMethodUpdate.Id = id

	c.paymentMethodService.UpdatePaymentMethod(paymentMethodUpdate)

	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusCreated, response)

}

func (c *paymentMethodController) FindPaymentMethodAll(ctx *gin.Context) {
	paymentMethodResponse := c.paymentMethodService.FindAll()
	response := helper.BuildResponse(true, "!OK", paymentMethodResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *paymentMethodController) FindPaymentMethodById(ctx *gin.Context) {
	paymentMethodId := ctx.Param("paymentMethodId")
	id, err := strconv.Atoi(paymentMethodId)
	helper.ErrorPanic(err)

	paymentMethodResponse := c.paymentMethodService.FindById(id)

	response := helper.BuildResponse(true, "!OK", paymentMethodResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *paymentMethodController) DeletePaymentMethod(ctx *gin.Context) {
	paymentMethodId := ctx.Param("paymentMethodId")
	id, err := strconv.Atoi(paymentMethodId)
	helper.ErrorPanic(err)

	c.paymentMethodService.Delete(id)

	response := helper.BuildResponse(true, "!OK", nil)
	ctx.JSON(http.StatusOK, response)
}
