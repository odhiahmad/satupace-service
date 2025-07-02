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
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	Delete(ctx *gin.Context)
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

func (c *paymentMethodController) Create(ctx *gin.Context) {
	var input request.PaymentMethodCreate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid",
			"INVALID_REQUEST",
			"body",
			err.Error(),
			nil,
		))
		return
	}

	res, err := c.paymentMethodService.CreatePaymentMethod(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal membuat metode pembayaran",
			"CREATE_FAILED",
			"internal",
			err.Error(),
			nil,
		))
		return
	}

	ctx.JSON(http.StatusCreated, helper.BuildResponse(true, "Berhasil membuat metode pembayaran", res))
}

func (c *paymentMethodController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter id wajib diisi",
			"MISSING_ID",
			"id",
			"Parameter 'id' tidak ditemukan dalam path",
			nil,
		))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid",
			"INVALID_ID",
			"id",
			err.Error(),
			nil,
		))
		return
	}

	var input request.PaymentMethodUpdate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid",
			"INVALID_REQUEST",
			"body",
			err.Error(),
			nil,
		))
		return
	}

	res, err := c.paymentMethodService.UpdatePaymentMethod(id, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengubah metode pembayaran",
			"UPDATE_FAILED",
			"internal",
			err.Error(),
			nil,
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengubah metode pembayaran", res))
}

func (c *paymentMethodController) FindById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid",
			"INVALID_ID",
			"id",
			err.Error(),
			nil,
		))
		return
	}

	res, err := c.paymentMethodService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, helper.BuildErrorResponse(
			"Metode pembayaran tidak ditemukan",
			"NOT_FOUND",
			"id",
			err.Error(),
			nil,
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil metode pembayaran", res))
}

func (c *paymentMethodController) FindAll(ctx *gin.Context) {
	res, err := c.paymentMethodService.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil semua metode pembayaran",
			"FETCH_FAILED",
			"internal",
			err.Error(),
			nil,
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil semua metode pembayaran", res))
}

func (c *paymentMethodController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid",
			"INVALID_ID",
			"id",
			err.Error(),
			nil,
		))
		return
	}

	err = c.paymentMethodService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal menghapus metode pembayaran",
			"DELETE_FAILED",
			"internal",
			err.Error(),
			nil,
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil menghapus metode pembayaran", nil))
}
