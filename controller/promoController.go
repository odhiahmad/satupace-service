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

type PromoController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindWithPagination(ctx *gin.Context)
	SetIsActive(ctx *gin.Context)
}

type promoController struct {
	promoService service.PromoService
	jwtService   service.JWTService
}

func NewPromoController(promoService service.PromoService, jwtService service.JWTService) PromoController {
	return &promoController{promoService: promoService, jwtService: jwtService}
}

func (c *promoController) Create(ctx *gin.Context) {
	var input request.PromoCreate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid",
			"BAD_REQUEST",
			"body",
			err.Error(),
			nil,
		))
		return
	}

	res, err := c.promoService.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal membuat promo",
			"CREATE_FAILED",
			"promo",
			err.Error(),
			nil,
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil membuat promo", res))
}

func (c *promoController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter id wajib diisi",
			"MISSING_ID",
			"id",
			"id diperlukan",
			nil,
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
			nil,
		))
		return
	}

	var input request.PromoUpdate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid",
			"BAD_REQUEST",
			"body",
			err.Error(),
			nil,
		))
		return
	}

	res, err := c.promoService.Update(id, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengubah promo",
			"UPDATE_FAILED",
			"promo",
			err.Error(),
			nil,
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengubah promo", res))
}

func (c *promoController) Delete(ctx *gin.Context) {
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

	err = c.promoService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal menghapus promo",
			"DELETE_FAILED",
			"promo",
			err.Error(),
			nil,
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil menghapus promo", nil))
}

func (c *promoController) FindById(ctx *gin.Context) {
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

	res, err := c.promoService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, helper.BuildErrorResponse(
			"Promo tidak ditemukan",
			"NOT_FOUND",
			"id",
			err.Error(),
			nil,
		))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil promo", res))
}

func (c *promoController) FindWithPagination(ctx *gin.Context) {
	businessIDStr := ctx.Query("business_id")
	if businessIDStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter business_id wajib diisi",
			"MISSING_BUSINESS_ID",
			"business_id",
			"business_id diperlukan",
			nil,
		))
		return
	}

	businessID, err := strconv.Atoi(businessIDStr)
	if err != nil || businessID <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter business_id tidak valid",
			"INVALID_BUSINESS_ID",
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
			"INVALID_LIMIT",
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

	promos, total, err := c.promoService.FindWithPagination(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data promo",
			"FETCH_FAILED",
			"promo",
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
		"Data promo berhasil diambil",
		promos,
		paginationMeta,
	))
}

func (c *promoController) SetIsActive(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid",
			"INVALID_ID",
			"id",
			err.Error(),
			nil,
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

	err = c.promoService.SetIsActive(id, input.IsActive)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengubah status promo",
			"UPDATE_STATUS_FAILED",
			"internal",
			err.Error(),
			nil,
		))
		return
	}

	status := "dinonaktifkan"
	if input.IsActive {
		status = "diaktifkan"
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Promo berhasil "+status, nil))
}
