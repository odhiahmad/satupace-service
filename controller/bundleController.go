package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type BundleController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	FindById(ctx *gin.Context)
	Delete(ctx *gin.Context)
	SetIsActive(ctx *gin.Context)
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
		res := helper.BuildErrorResponse(
			"Gagal memproses data",
			"VALIDATION_ERROR",
			"request_body",
			err.Error(),
			nil,
		)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.bundleService.CreateBundle(req); err != nil {
		res := helper.BuildErrorResponse(
			"Gagal membuat bundle",
			"BUNDLE_CREATE_FAILED",
			"",
			err.Error(),
			nil,
		)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Bundle berhasil dibuat", helper.EmptyObj{})
	ctx.JSON(http.StatusCreated, res)
}

func (c *bundleController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse(
			"ID tidak valid",
			"INVALID_ID",
			"id",
			err.Error(),
			nil,
		)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var req request.BundleUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse(
			"Gagal memproses data",
			"VALIDATION_ERROR",
			"request_body",
			err.Error(),
			nil,
		)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.bundleService.UpdateBundle(id, req); err != nil {
		res := helper.BuildErrorResponse(
			"Gagal memperbarui bundle",
			"BUNDLE_UPDATE_FAILED",
			"",
			err.Error(),
			nil,
		)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Bundle berhasil diperbarui", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (c *bundleController) FindById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse(
			"ID tidak valid",
			"INVALID_ID",
			"id",
			err.Error(),
			nil,
		)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	bundle, err := c.bundleService.FindById(id)
	if err != nil {
		res := helper.BuildErrorResponse(
			"Bundle tidak ditemukan",
			"BUNDLE_NOT_FOUND",
			"id",
			err.Error(),
			nil,
		)
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := helper.BuildResponse(true, "Bundle berhasil ditemukan", bundle)
	ctx.JSON(http.StatusOK, res)
}

func (c *bundleController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse(
			"ID tidak valid",
			"INVALID_ID",
			"id",
			err.Error(),
			nil,
		)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.bundleService.Delete(id); err != nil {
		res := helper.BuildErrorResponse(
			"Gagal menghapus bundle",
			"BUNDLE_DELETE_FAILED",
			"",
			err.Error(),
			nil,
		)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Bundle berhasil dihapus", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (c *bundleController) SetIsActive(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter id wajib diisi",
			"MISSING_ID",
			"id",
			"Parameter 'id' tidak ditemukan dalam path",
			helper.EmptyObj{},
		))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter id tidak valid",
			"INVALID_ID",
			"id",
			err.Error(),
			helper.EmptyObj{},
		))
		return
	}

	var input request.IsActive
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

	err = c.bundleService.SetIsActive(id, input.IsActive)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengubah status diskon",
			"UPDATE_STATUS_FAILED",
			"internal",
			err.Error(),
			nil,
		))
		return
	}

	statusMsg := "dinonaktifkan"
	if input.IsActive {
		statusMsg = "diaktifkan"
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Diskon berhasil "+statusMsg, nil))
}

func (c *bundleController) FindWithPagination(ctx *gin.Context) {
	businessIDStr := ctx.Query("business_id")
	if businessIDStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter business_id wajib diisi",
			"REQUIRED_QUERY_PARAM",
			"business_id",
			"Parameter business_id tidak boleh kosong",
			nil,
		))
		return
	}

	businessID, err := strconv.Atoi(businessIDStr)
	if err != nil || businessID <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter business_id tidak valid",
			"INVALID_QUERY_PARAM",
			"business_id",
			err.Error(),
			nil,
		))
		return
	}

	limitStr := ctx.DefaultQuery("limit", "10")
	sortBy := ctx.DefaultQuery("sortBy", "created_at")
	orderBy := ctx.DefaultQuery("orderBy", "desc")
	search := ctx.DefaultQuery("search", "")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter limit tidak valid",
			"INVALID_QUERY_PARAM",
			"limit",
			err.Error(),
			nil,
		))
		return
	}

	pagination := request.Pagination{
		Limit:   limit,
		SortBy:  sortBy,
		OrderBy: orderBy,
		Search:  search,
	}

	bundles, total, err := c.bundleService.FindWithPagination(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data bundle",
			"BUNDLE_FETCH_FAILED",
			"",
			err.Error(),
			nil,
		))
		return
	}

	paginationMeta := response.PaginatedResponse{
		Page:      pagination.Page,
		Limit:     pagination.Limit,
		Total:     total,
		OrderBy:   pagination.SortBy,
		SortOrder: pagination.OrderBy,
	}

	ctx.JSON(http.StatusOK, helper.BuildResponsePagination(
		true,
		"Data bundle berhasil diambil",
		bundles,
		paginationMeta,
	))
}
