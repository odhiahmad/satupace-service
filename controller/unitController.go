package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type UnitController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindWithPagination(ctx *gin.Context)
	FindWithPaginationCursor(ctx *gin.Context)
}

type unitController struct {
	service    service.UnitService
	jwtService service.JWTService
}

func NewUnitController(s service.UnitService, jwtService service.JWTService) UnitController {
	return &unitController{service: s, jwtService: jwtService}
}

func (c *unitController) Create(ctx *gin.Context) {
	businessIdStr := ctx.MustGet("business_id").(string)
	businessId, err := uuid.Parse(businessIdStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid business_id UUID"})
		return
	}

	var input request.UnitRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid", "BAD_REQUEST", "body", err.Error(), helper.EmptyObj{}))
		return
	}
	input.BusinessId = businessId

	res, err := c.service.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal membuat satuan produk", "INTERNAL_ERROR", "product_unit", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusCreated, helper.BuildResponse(true, "Berhasil membuat satuan produk", res))
}

// Update function: UnitController
func (c *unitController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	businessIdStr := ctx.MustGet("business_id").(string)
	businessId, err := uuid.Parse(businessIdStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid business_id UUID"})
		return
	}

	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter id wajib diisi", "missing_parameter", "id", "parameter id kosong", nil))
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter id tidak valid", "invalid_parameter", "id", err.Error(), nil))
		return
	}

	var input request.UnitRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "BAD_REQUEST", "body", err.Error(), helper.EmptyObj{}))
		return
	}
	input.BusinessId = businessId

	res, err := c.service.Update(id, input) // <-- Pastikan service menerima uuid.UUID
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengubah satuan produk", "INTERNAL_ERROR", "product_unit", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengubah satuan produk", res))
}

func (c *unitController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", "BAD_REQUEST", "id", err.Error(), helper.EmptyObj{}))
		return
	}

	err = c.service.Delete(id) // <-- Pastikan service menerima uuid.UUID
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal menghapus satuan produk", "INTERNAL_ERROR", "product_unit", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil menghapus satuan produk", helper.EmptyObj{}))
}

func (c *unitController) FindById(ctx *gin.Context) {
	unitIdStr := ctx.Param("id")
	unitId, err := uuid.Parse(unitIdStr)
	if err != nil {
		response := helper.BuildErrorResponse("Parameter id tidak valid", "invalid_parameter", "id", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	unitResponse := c.service.FindById(unitId) // <-- Pastikan service menerima uuid.UUID

	response := helper.BuildResponse(true, "Berhasil mengambil data unit", unitResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *unitController) FindWithPagination(ctx *gin.Context) {
	businessIdStr := ctx.MustGet("business_id").(string)
	businessID, err := uuid.Parse(businessIdStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid business_id UUID"})
		return
	}
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	orderBy := ctx.DefaultQuery("order_by", "desc")
	search := ctx.DefaultQuery("search", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter page tidak valid", "BAD_REQUEST", "page", err.Error(), helper.EmptyObj{}))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter limit tidak valid", "BAD_REQUEST", "limit", err.Error(), helper.EmptyObj{}))
		return
	}

	pagination := request.Pagination{
		Page:    page,
		Limit:   limit,
		SortBy:  sortBy,
		OrderBy: orderBy,
		Search:  search,
	}

	units, total, err := c.service.FindWithPagination(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data unit produk", "INTERNAL_ERROR", "unit", err.Error(), helper.EmptyObj{}))
		return
	}

	paginationMeta := response.PaginatedResponse{
		Page:      page,
		Limit:     limit,
		Total:     total,
		OrderBy:   sortBy,
		SortOrder: orderBy,
	}

	ctx.JSON(http.StatusOK, helper.BuildResponsePagination(
		true,
		"Berhasil mengambil data unit produk",
		units,
		paginationMeta,
	))
}

func (c *unitController) FindWithPaginationCursor(ctx *gin.Context) {
	businessIdStr := ctx.MustGet("business_id").(string)
	businessID, err := uuid.Parse(businessIdStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid business_id UUID"})
		return
	}
	limitStr := ctx.DefaultQuery("limit", "10")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	orderBy := ctx.DefaultQuery("order_by", "desc")
	search := ctx.DefaultQuery("search", "")
	cursor := ctx.DefaultQuery("cursor", "")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter limit tidak valid", "BAD_REQUEST", "limit", err.Error(), helper.EmptyObj{}))
		return
	}

	pagination := request.Pagination{
		Cursor:  cursor,
		Limit:   limit,
		SortBy:  sortBy,
		OrderBy: orderBy,
		Search:  search,
	}

	units, nextCursor, hasNext, err := c.service.FindWithPaginationCursor(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data brand", "internal_error", "brand", err.Error(), nil))
		return
	}

	paginationMeta := response.CursorPaginatedResponse{
		Limit:      limit,
		SortBy:     sortBy,
		OrderBy:    orderBy,
		NextCursor: nextCursor,
		HasNext:    hasNext,
	}

	ctx.JSON(http.StatusOK, helper.BuildResponseCursorPagination(
		true,
		"Data brand berhasil diambil",
		units,
		paginationMeta,
	))
}
