package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type ProductVariantController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindByProductId(ctx *gin.Context)
	SetActive(ctx *gin.Context)
	SetAvailable(ctx *gin.Context)
}

type productVariantController struct {
	Service    service.ProductVariantService
	JWTService service.JWTService
}

func NewProductVariantController(s service.ProductVariantService, jwt service.JWTService) ProductVariantController {
	return &productVariantController{Service: s, JWTService: jwt}
}

func (c *productVariantController) Create(ctx *gin.Context) {
	var req request.ProductVariantRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid", "BAD_REQUEST", "body", err.Error(), helper.EmptyObj{}))
		return
	}

	productIdStr := ctx.Param("id")
	productId, err := uuid.Parse(productIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter product_id tidak valid", "BAD_REQUEST", "product_id", err.Error(), helper.EmptyObj{}))
		return
	}

	variant, err := c.Service.Create(req, productId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal membuat variant produk", "INTERNAL_ERROR", "product_variant", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusCreated, helper.BuildResponse(true, "Variant berhasil dibuat", variant))
}

func (c *productVariantController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid", "BAD_REQUEST", "id", err.Error(), helper.EmptyObj{}))
		return
	}

	var req request.ProductVariantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid", "BAD_REQUEST", "body", err.Error(), helper.EmptyObj{}))
		return
	}
	variant, err := c.Service.Update(id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal memperbarui variant", "INTERNAL_ERROR", "product_variant", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Variant berhasil diperbarui", variant))
}

func (c *productVariantController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid", "BAD_REQUEST", "id", err.Error(), helper.EmptyObj{}))
		return
	}

	if err := c.Service.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal menghapus variant", "INTERNAL_ERROR", "product_variant", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Variant berhasil dihapus", helper.EmptyObj{}))
}

func (c *productVariantController) FindById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid", "BAD_REQUEST", "id", err.Error(), helper.EmptyObj{}))
		return
	}

	variant, err := c.Service.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, helper.BuildErrorResponse(
			"Variant tidak ditemukan", "NOT_FOUND", "product_variant", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Variant ditemukan", variant))
}

func (c *productVariantController) FindByProductId(ctx *gin.Context) {
	productIdStr := ctx.Param("productId")
	productId, err := uuid.Parse(productIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Product ID tidak valid", "BAD_REQUEST", "product_id", err.Error(), helper.EmptyObj{}))
		return
	}

	variants, err := c.Service.FindByProductId(productId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data variant", "INTERNAL_ERROR", "product_variant", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Data variant ditemukan", variants))
}

func (c *productVariantController) SetActive(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid", "BAD_REQUEST", "id", err.Error(), helper.EmptyObj{}))
		return
	}

	var body struct {
		IsActive bool `json:"is_active"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid", "BAD_REQUEST", "body", err.Error(), helper.EmptyObj{}))
		return
	}

	if err := c.Service.SetActive(id, body.IsActive); err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengubah status aktif", "INTERNAL_ERROR", "product_variant", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Status aktif variant diperbarui", helper.EmptyObj{}))
}

func (c *productVariantController) SetAvailable(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid", "BAD_REQUEST", "id", err.Error(), helper.EmptyObj{}))
		return
	}

	var body struct {
		IsAvailable bool `json:"is_available"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid", "BAD_REQUEST", "body", err.Error(), helper.EmptyObj{}))
		return
	}

	if err := c.Service.SetAvailable(id, body.IsAvailable); err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengubah status ketersediaan", "INTERNAL_ERROR", "product_variant", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Status ketersediaan variant diperbarui", helper.EmptyObj{}))
}
