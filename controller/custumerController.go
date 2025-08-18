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

type CustomerController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindWithPagination(ctx *gin.Context)
	FindWithPaginationCursor(ctx *gin.Context)
}

type customerController struct {
	customerService service.CustomerService
	jwtService      service.JWTService
}

func NewCustomerController(customerService service.CustomerService, jwtService service.JWTService) CustomerController {
	return &customerController{customerService: customerService, jwtService: jwtService}
}

func (c *customerController) Create(ctx *gin.Context) {
	businessIdStr := ctx.MustGet("business_id").(string)
	businessId, err := uuid.Parse(businessIdStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid business_id UUID"})
		return
	}
	var input request.CustomerRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "bad_request", "body", err.Error(), nil))
		return
	}

	input.BusinessId = businessId

	res, err := c.customerService.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal membuat customer", "internal_error", "customer", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, helper.BuildResponse(true, "Berhasil membuat customer", res))
}

func (c *customerController) Update(ctx *gin.Context) {
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

	var input request.CustomerRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "bad_request", "body", err.Error(), nil))
		return
	}

	input.BusinessId = businessId

	res, err := c.customerService.Update(id, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengubah customer", "internal_error", "customer", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengubah customer", res))
}

func (c *customerController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", "invalid_parameter", "id", err.Error(), nil))
		return
	}

	err = c.customerService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal menghapus customer", "internal_error", "customer", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil menghapus customer", nil))
}

func (c *customerController) FindById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter id tidak valid", "invalid_parameter", "id", err.Error(), nil))
		return
	}

	customerResponse := c.customerService.FindById(id)
	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil data customer", customerResponse))
}

func (c *customerController) FindWithPagination(ctx *gin.Context) {
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

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter limit tidak valid", "invalid_parameter", "limit", err.Error(), nil))
		return
	}

	pagination := request.Pagination{
		Limit:   limit,
		SortBy:  sortBy,
		OrderBy: orderBy,
		Search:  search,
	}

	customeres, total, err := c.customerService.FindWithPagination(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengambil data customer", "internal_error", "customer", err.Error(), nil))
		return
	}

	paginationMeta := response.PaginatedResponse{
		Page:      1,
		Limit:     pagination.Limit,
		Total:     total,
		OrderBy:   pagination.SortBy,
		SortOrder: pagination.OrderBy,
	}

	ctx.JSON(http.StatusOK, helper.BuildResponsePagination(
		true,
		"Data customer berhasil diambil",
		customeres,
		paginationMeta,
	))
}

func (c *customerController) FindWithPaginationCursor(ctx *gin.Context) {
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
			"Parameter limit tidak valid", "invalid_parameter", "limit", err.Error(), nil))
		return
	}

	pagination := request.Pagination{
		Cursor:  cursor,
		Limit:   limit,
		SortBy:  sortBy,
		OrderBy: orderBy,
		Search:  search,
	}

	customers, nextCursor, hasNext, err := c.customerService.FindWithPaginationCursor(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data customer", "internal_error", "customer", err.Error(), nil))
		return
	}

	paginationMeta := response.CursorPaginatedResponse{
		Limit:      limit,
		SortBy:     sortBy,
		OrderBy:    orderBy,
		NextCursor: nextCursor,
		HasNext:    hasNext, // konversi bool ke string: "true"/"false"
	}

	ctx.JSON(http.StatusOK, helper.BuildResponseCursorPagination(
		true,
		"Data customer berhasil diambil",
		customers,
		paginationMeta,
	))
}
