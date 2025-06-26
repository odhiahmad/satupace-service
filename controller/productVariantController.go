package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type ProductVariantController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	DeleteByProductId(ctx *gin.Context)
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
	var req request.ProductVariantCreate

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Gagal memproses permintaan", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	productIdStr := ctx.Param("id")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		res := helper.BuildErrorResponse("Parameter product_id tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	variant, err := c.Service.Create(req, productId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal membuat variant produk", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := helper.BuildResponse(true, "Variant berhasil dibuat", variant)
	ctx.JSON(http.StatusCreated, res)
}

func (c *productVariantController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse("ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var req request.ProductVariantUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Gagal memproses permintaan", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.Service.Update(id, req); err != nil {
		res := helper.BuildErrorResponse("Gagal memperbarui variant", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := helper.BuildResponse(true, "Variant berhasil diperbarui", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (c *productVariantController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse("ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.Service.Delete(id); err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus variant", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := helper.BuildResponse(true, "Variant berhasil dihapus", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (c *productVariantController) DeleteByProductId(ctx *gin.Context) {
	productIdStr := ctx.Param("productId")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		res := helper.BuildErrorResponse("Product ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.Service.DeleteByProductId(productId); err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus seluruh variant", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := helper.BuildResponse(true, "Semua variant berhasil dihapus", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (c *productVariantController) FindById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse("ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	variant, err := c.Service.FindById(id)
	if err != nil {
		res := helper.BuildErrorResponse("Variant tidak ditemukan", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := helper.BuildResponse(true, "Variant ditemukan", variant)
	ctx.JSON(http.StatusOK, res)
}

func (c *productVariantController) FindByProductId(ctx *gin.Context) {
	productIdStr := ctx.Param("productId")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		res := helper.BuildErrorResponse("Product ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	variants, err := c.Service.FindByProductId(productId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil data variant", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := helper.BuildResponse(true, "Data variant ditemukan", variants)
	ctx.JSON(http.StatusOK, res)
}

func (c *productVariantController) SetActive(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product variant ID"})
		return
	}

	var body struct {
		IsActive bool `json:"is_active"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.Service.SetActive(id, body.IsActive)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product variant active status updated"})
}

func (c *productVariantController) SetAvailable(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product variant ID"})
		return
	}

	var body struct {
		IsAvailable bool `json:"is_available"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.Service.SetAvailable(id, body.IsAvailable)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product variant availability status updated"})
}
