package controller

import (
	"net/http"
	"strconv"

	"loka-kasir/data/request"
	"loka-kasir/data/response"
	"loka-kasir/helper"
	"loka-kasir/service"

	"github.com/gin-gonic/gin"
)

type OrderTypeController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindWithPagination(ctx *gin.Context)
	FindWithPaginationCursor(ctx *gin.Context)
}

type orderTypeController struct {
	orderTypeService service.OrderTypeService
	jwtService       service.JWTService
}

func NewOrderTypeController(orderTypeService service.OrderTypeService, jwtService service.JWTService) OrderTypeController {
	return &orderTypeController{orderTypeService: orderTypeService, jwtService: jwtService}
}

func (c *orderTypeController) Create(ctx *gin.Context) {
	var input request.OrderTypeRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "bad_request", "body", err.Error(), nil))
		return
	}

	res, err := c.orderTypeService.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal membuat orderType", "internal_error", "orderType", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, helper.BuildResponse(true, "Berhasil membuat orderType", res))
}

func (c *orderTypeController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter id wajib diisi", "missing_parameter", "id", "parameter id kosong", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter id tidak valid", "invalid_parameter", "id", err.Error(), nil))
		return
	}

	var input request.OrderTypeRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "bad_request", "body", err.Error(), nil))
		return
	}

	res, err := c.orderTypeService.Update(id, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengubah orderType", "internal_error", "orderType", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengubah orderType", res))
}

func (c *orderTypeController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", "invalid_parameter", "id", err.Error(), nil))
		return
	}

	err = c.orderTypeService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal menghapus orderType", "internal_error", "orderType", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil menghapus orderType", nil))
}

func (c *orderTypeController) FindById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter id tidak valid", "invalid_parameter", "id", err.Error(), nil))
		return
	}

	orderTypeResponse := c.orderTypeService.FindById(id)
	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil data orderType", orderTypeResponse))
}

func (c *orderTypeController) FindWithPagination(ctx *gin.Context) {
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

	orderTypees, total, err := c.orderTypeService.FindWithPagination(pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengambil data orderType", "internal_error", "orderType", err.Error(), nil))
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
		"Data orderType berhasil diambil",
		orderTypees,
		paginationMeta,
	))
}

func (c *orderTypeController) FindWithPaginationCursor(ctx *gin.Context) {
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

	orderTypes, nextCursor, hasNext, err := c.orderTypeService.FindWithPaginationCursor(pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data orderType", "internal_error", "orderType", err.Error(), nil))
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
		"Data orderType berhasil diambil",
		orderTypes,
		paginationMeta,
	))
}
