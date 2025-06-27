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

type TaxController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindWithPagination(ctx *gin.Context)
}

type taxController struct {
	taxService service.TaxService
	jwtService service.JWTService
}

func NewTaxController(taxService service.TaxService, jwtService service.JWTService) TaxController {
	return &taxController{taxService: taxService, jwtService: jwtService}
}

func (c *taxController) Create(ctx *gin.Context) {
	var input request.TaxCreate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", err.Error(), nil))
		return
	}

	res, err := c.taxService.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal membuat tax", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil membuat tax", res))
}

func (c *taxController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter id wajib diisi", "missing id", helper.EmptyObj{}))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter id tidak valid", err.Error(), helper.EmptyObj{}))
		return
	}

	var input request.TaxUpdate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", err.Error(), nil))
		return
	}

	res, err := c.taxService.Update(id, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengubah tax", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengubah tax", res))
}

func (c *taxController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", err.Error(), nil))
		return
	}

	err = c.taxService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal menghapus tax", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil menghapus tax", nil))
}

func (c *taxController) FindById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", err.Error(), nil))
		return
	}

	res, err := c.taxService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, helper.BuildErrorResponse("Tax tidak ditemukan", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil tax", res))
}

func (c *taxController) FindWithPagination(ctx *gin.Context) {
	businessIDStr := ctx.Query("business_id")
	if businessIDStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter business_id wajib diisi", "missing business_id", helper.EmptyObj{}))
		return
	}

	businessID, err := strconv.Atoi(businessIDStr)
	if err != nil || businessID <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter business_id tidak valid", err.Error(), helper.EmptyObj{}))
		return
	}

	// Ambil dan parsing query parameter pagination
	limitStr := ctx.DefaultQuery("limit", "10")
	sortBy := ctx.DefaultQuery("sortBy", "created_at")
	orderBy := ctx.DefaultQuery("orderBy", "desc")
	search := ctx.DefaultQuery("search", "")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter limit tidak valid", err.Error(), helper.EmptyObj{}))
		return
	}

	// Susun struct pagination
	pagination := request.Pagination{
		Limit:   limit,
		SortBy:  sortBy,
		OrderBy: orderBy,
		Search:  search,
	}

	// Ambil data dari service
	taxes, total, err := c.taxService.FindWithPagination(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data tax", err.Error(), helper.EmptyObj{}))
		return
	}

	// Susun metadata pagination
	paginationMeta := response.PaginatedResponse{
		Page:      1, // Default 1 jika tidak ada page
		Limit:     pagination.Limit,
		Total:     total,
		OrderBy:   pagination.SortBy,
		SortOrder: pagination.OrderBy,
	}

	// Kirim response
	ctx.JSON(http.StatusOK, helper.BuildResponsePagination(
		true,
		"Data tax berhasil diambil",
		taxes,
		paginationMeta,
	))
}
