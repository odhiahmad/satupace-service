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
	Delete(ctx *gin.Context)
	FindWithPagination(ctx *gin.Context)
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
		res := helper.BuildErrorResponse("Gagal memproses data", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.bundleService.CreateBundle(req); err != nil {
		res := helper.BuildErrorResponse("Gagal membuat bundle", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Bundle berhasil dibuat", helper.EmptyObj{})
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
		res := helper.BuildErrorResponse("Gagal memproses data", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.bundleService.UpdateBundle(id, req); err != nil {
		res := helper.BuildErrorResponse("Gagal memperbarui bundle", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Bundle berhasil diperbarui", helper.EmptyObj{})
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
		res := helper.BuildErrorResponse("Bundle tidak ditemukan", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := helper.BuildResponse(true, "Bundle berhasil ditemukan", bundle)
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
		res := helper.BuildErrorResponse("Gagal menghapus bundle", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Bundle berhasil dihapus", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

// Find All with Pagination
func (c *bundleController) FindWithPagination(ctx *gin.Context) {
	businessIDStr := ctx.Query("business_id")
	if businessIDStr == "" {
		res := helper.BuildErrorResponse("Parameter business_id wajib diisi", "missing business_id", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	businessID, err := strconv.Atoi(businessIDStr)
	if err != nil || businessID <= 0 {
		res := helper.BuildErrorResponse("Parameter business_id tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Ambil query parameter lainnya
	limitStr := ctx.DefaultQuery("limit", "10")
	sortBy := ctx.DefaultQuery("sortBy", "id")
	orderBy := ctx.DefaultQuery("orderBy", "asc")
	search := ctx.DefaultQuery("search", "")
	before := ctx.Query("before")
	after := ctx.Query("after")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		res := helper.BuildErrorResponse("Parameter limit tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	pagination := request.Pagination{
		Limit:   limit,
		SortBy:  sortBy,
		OrderBy: orderBy,
		Search:  search,
		Before:  before,
		After:   after,
	}

	bundles, total, err := c.bundleService.FindWithPagination(businessID, pagination)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil data bundle", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// Bisa disesuaikan jika kamu ingin kirim next/prev cursor juga
	responseData := gin.H{
		"total":   total,
		"limit":   limit,
		"results": bundles,
	}

	res := helper.BuildResponse(true, "Bundle berhasil diambil", responseData)
	ctx.JSON(http.StatusOK, res)
}
