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

type DiscountController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindWithPagination(ctx *gin.Context)
	SetIsActive(ctx *gin.Context)
}

type discountController struct {
	discountService service.DiscountService
	jwtService      service.JWTService
}

func NewDiscountController(discountService service.DiscountService, jwtService service.JWTService) DiscountController {
	return &discountController{discountService: discountService, jwtService: jwtService}
}

func (c *discountController) Create(ctx *gin.Context) {
	businessId := ctx.MustGet("business_id").(int)
	var input request.DiscountRequest
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

	input.BusinessId = businessId

	res, err := c.discountService.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal membuat diskon",
			"CREATE_FAILED",
			"internal",
			err.Error(),
			nil,
		))
		return
	}

	ctx.JSON(http.StatusCreated, helper.BuildResponse(true, "Berhasil membuat diskon", res))
}

func (c *discountController) Update(ctx *gin.Context) {
	businessId := ctx.MustGet("business_id").(int)
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

	var input request.DiscountRequest
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

	input.BusinessId = businessId
	res, err := c.discountService.Update(id, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengubah diskon",
			"UPDATE_FAILED",
			"internal",
			err.Error(),
			nil,
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengubah diskon", res))
}

func (c *discountController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid",
			"INVALID_ID",
			"id",
			err.Error(),
			nil,
		))
		return
	}

	err = c.discountService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal menghapus diskon",
			"DELETE_FAILED",
			"internal",
			err.Error(),
			nil,
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil menghapus diskon", nil))
}

func (c *discountController) FindById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid",
			"INVALID_ID",
			"id",
			err.Error(),
			nil,
		))
		return
	}

	res, err := c.discountService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, helper.BuildErrorResponse(
			"Diskon tidak ditemukan",
			"NOT_FOUND",
			"id",
			err.Error(),
			nil,
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil diskon", res))
}

func (c *discountController) FindWithPagination(ctx *gin.Context) {
	businessID := ctx.MustGet("business_id").(int)
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	orderBy := ctx.DefaultQuery("order_by", "desc")
	search := ctx.DefaultQuery("search", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter page tidak valid",
			"INVALID_PAGE",
			"page",
			err.Error(),
			helper.EmptyObj{},
		))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter limit tidak valid",
			"INVALID_LIMIT",
			"limit",
			err.Error(),
			helper.EmptyObj{},
		))
		return
	}

	pagination := request.Pagination{
		Page:    page,
		Limit:   limit,
		SortBy:  sortBy,
		OrderBy: orderBy,
		Search:  search,
	}

	discounts, total, err := c.discountService.FindWithPagination(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data diskon",
			"FETCH_FAILED",
			"internal",
			err.Error(),
			helper.EmptyObj{},
		))
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
		"Berhasil mengambil data diskon",
		discounts,
		paginationMeta,
	))
}

func (c *discountController) SetIsActive(ctx *gin.Context) {
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

	err = c.discountService.SetIsActive(id, input.IsActive)
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
