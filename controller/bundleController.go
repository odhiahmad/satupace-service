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
	SetIsAvailable(ctx *gin.Context)
	FindWithPagination(ctx *gin.Context)
	FindWithPaginationCursor(ctx *gin.Context)
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
	businessId := ctx.MustGet("business_id").(int)
	var input request.BundleRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Gagal bind data bundle", "INVALID_JSON", "body", err.Error(), helper.EmptyObj{}))
		return
	}

	input.BusinessId = businessId
	res, err := c.bundleService.CreateBundle(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal membuat bundle", "internal_error", "brand", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil membuat bundle", res))
}

func (c *bundleController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	businessId := ctx.MustGet("business_id").(int)

	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter id wajib diisi", "missing_parameter", "id", "parameter id kosong", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter id tidak valid", "invalid_parameter", "id", err.Error(), nil))
		return
	}

	var input request.BundleRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "bad_request", "body", err.Error(), nil))
		return
	}

	input.BusinessId = businessId

	res, err := c.bundleService.UpdateBundle(id, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengubah bundle", "internal_error", "bundle", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Bundle berhasil diperbaru", res))
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

func (c *bundleController) SetIsAvailable(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid", "INVALID_ID", "id", err.Error(), helper.EmptyObj{}))
		return
	}

	var body struct {
		IsAvailable bool `json:"is_available"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Gagal bind body", "INVALID_JSON", "body", err.Error(), helper.EmptyObj{}))
		return
	}

	if err := c.bundleService.SetIsAvailable(id, body.IsAvailable); err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengubah status ketersediaan bundle", "UPDATE_ERROR", "is_available", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Status ketersediaan bundle berhasil diubah", helper.EmptyObj{}))
}

func (c *bundleController) FindWithPagination(ctx *gin.Context) {
	businessID := ctx.MustGet("business_id").(int)
	limitStr := ctx.DefaultQuery("limit", "10")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	orderBy := ctx.DefaultQuery("order_by", "desc")
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

func (c *bundleController) FindWithPaginationCursor(ctx *gin.Context) {
	businessID := ctx.MustGet("business_id").(int)
	limitStr := ctx.DefaultQuery("limit", "10")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	orderBy := ctx.DefaultQuery("order_by", "desc")
	search := ctx.DefaultQuery("search", "")
	cursor := ctx.DefaultQuery("cursor", "")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter limit tidak valid", "INVALID_QUERY_PARAM", "limit", err.Error(), nil))
		return
	}

	pagination := request.Pagination{
		Cursor:  cursor,
		Limit:   limit,
		SortBy:  sortBy,
		OrderBy: orderBy,
		Search:  search,
	}

	bundles, nextCursor, err := c.bundleService.FindWithPaginationCursor(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data bundle", "BUNDLE_FETCH_FAILED", "bundle", err.Error(), nil))
		return
	}

	paginationMeta := response.CursorPaginatedResponse{
		Limit:      limit,
		SortBy:     sortBy,
		OrderBy:    orderBy,
		NextCursor: nextCursor,
	}

	ctx.JSON(http.StatusOK, helper.BuildResponseCursorPagination(
		true,
		"Data bundle berhasil diambil",
		bundles,
		paginationMeta,
	))
}
