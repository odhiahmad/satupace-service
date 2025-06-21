package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type BundleController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindByBusinessId(ctx *gin.Context)
}

type bundleController struct {
	bundleService service.BundleService
	jwtService    service.JWTService
}

func NewBundleController(bundleService service.BundleService, jwtService service.JWTService) BundleController {
	return &bundleController{
		bundleService: bundleService,
		jwtService:    jwtService,
	}
}

// Create Bundle
func (c *bundleController) Create(ctx *gin.Context) {
	var req request.BundleCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Gagal bind data", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.bundleService.CreateBundle(req); err != nil {
		res := helper.BuildErrorResponse("Gagal membuat produk", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Produk berhasil dibuat", helper.EmptyObj{})
	ctx.JSON(http.StatusCreated, res)
}

// Update Bundle
func (c *bundleController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse("ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var req request.BundleUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Gagal bind data", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.bundleService.UpdateBundle(id, req); err != nil {
		res := helper.BuildErrorResponse("Gagal mengupdate produk", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Produk berhasil diperbarui", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

// Find By ID
func (c *bundleController) FindById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse("ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	bundle, err := c.bundleService.FindById(id)
	if err != nil {
		res := helper.BuildErrorResponse("Produk tidak ditemukan", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := helper.BuildResponse(true, "Produk berhasil ditemukan", bundle)
	ctx.JSON(http.StatusOK, res)
}

// Find All
func (c *bundleController) FindAll(ctx *gin.Context) {
	bundles, err := c.bundleService.FindAll()
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil data produk", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := helper.BuildResponse(true, "Semua produk berhasil diambil", bundles)
	ctx.JSON(http.StatusOK, res)
}

// Delete Bundle
func (c *bundleController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse("ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.bundleService.Delete(id); err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus produk", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Produk berhasil dihapus", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (c *bundleController) FindByBusinessId(ctx *gin.Context) {
	businessIdStr := ctx.Param("business_id")
	businessId, err := strconv.Atoi(businessIdStr)
	if err != nil {
		res := helper.BuildErrorResponse("Business ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	bundles, err := c.bundleService.FindByBusinessId(businessId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil data produk berdasarkan bisnis", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := helper.BuildResponse(true, "Produk berdasarkan bisnis berhasil diambil", bundles)
	ctx.JSON(http.StatusOK, res)
}
