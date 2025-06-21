package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type ProductCategoryController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	FindByBusinessId(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type productCategoryController struct {
	productCategoryService service.ProductCategoryService
	jwtService             service.JWTService
}

func NewProductCategoryController(productCategoryService service.ProductCategoryService, jwtService service.JWTService) ProductCategoryController {
	return &productCategoryController{
		productCategoryService: productCategoryService,
		jwtService:             jwtService,
	}
}

// Create
func (c *productCategoryController) Create(ctx *gin.Context) {
	var req request.ProductCategoryCreate
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal bind data", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Validasi BusinessId tidak boleh kosong (0)
	if req.BusinessId == 0 {
		res := helper.BuildErrorResponse("Gagal membuat kategori", "BusinessId tidak boleh kosong", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Lanjut ke service
	err = c.productCategoryService.Create(req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal membuat kategori", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Berhasil membuat kategori produk", helper.EmptyObj{})
	ctx.JSON(http.StatusCreated, res)
}

// Update
func (c *productCategoryController) Update(ctx *gin.Context) {
	var req request.ProductCategoryUpdate
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal bind data", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err = c.productCategoryService.Update(req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal update kategori", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Berhasil update kategori produk", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

// FindById
func (c *productCategoryController) FindById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse("ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.productCategoryService.FindById(id)
	if err != nil {
		res := helper.BuildErrorResponse("Kategori tidak ditemukan", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := helper.BuildResponse(true, "Berhasil mengambil kategori produk", data)
	ctx.JSON(http.StatusOK, res)
}

// FindAll
func (c *productCategoryController) FindAll(ctx *gin.Context) {
	data, err := c.productCategoryService.FindAll()
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil data", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := helper.BuildResponse(true, "Berhasil mengambil semua kategori produk", data)
	ctx.JSON(http.StatusOK, res)
}

func (c *productCategoryController) FindByBusinessId(ctx *gin.Context) {
	businessIdStr := ctx.Param("business_id")
	businessId, err := strconv.Atoi(businessIdStr)
	if err != nil {
		res := helper.BuildErrorResponse("Business ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := c.productCategoryService.FindByBusinessId(businessId)
	if err != nil {
		res := helper.BuildErrorResponse("Kategori tidak ditemukan", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := helper.BuildResponse(true, "Berhasil mengambil kategori berdasarkan bisnis", data)
	ctx.JSON(http.StatusOK, res)
}

// Delete
func (c *productCategoryController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse("ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err = c.productCategoryService.Delete(id)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus kategori", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Berhasil menghapus kategori produk", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}
