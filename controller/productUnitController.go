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

type ProductUnitController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindWithPagination(ctx *gin.Context)
}

type productUnitController struct {
	service    service.ProductUnitService
	jwtService service.JWTService
}

func NewProductUnitController(s service.ProductUnitService, jwtService service.JWTService) ProductUnitController {
	return &productUnitController{service: s, jwtService: jwtService}
}

func (c *productUnitController) Create(ctx *gin.Context) {
	var input request.ProductUnitCreate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid", "BAD_REQUEST", "body", err.Error(), helper.EmptyObj{}))
		return
	}

	res, err := c.service.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal membuat satuan produk", "INTERNAL_ERROR", "product_unit", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil membuat satuan produk", res))
}

func (c *productUnitController) Update(ctx *gin.Context) {
	var input request.ProductUnitUpdate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Input tidak valid", "BAD_REQUEST", "body", err.Error(), helper.EmptyObj{}))
		return
	}

	res, err := c.service.Update(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengubah satuan produk", "INTERNAL_ERROR", "product_unit", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengubah satuan produk", res))
}

func (c *productUnitController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid", "BAD_REQUEST", "id", err.Error(), helper.EmptyObj{}))
		return
	}

	err = c.service.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal menghapus satuan produk", "INTERNAL_ERROR", "product_unit", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil menghapus satuan produk", helper.EmptyObj{}))
}

func (c *productUnitController) FindById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid", "BAD_REQUEST", "id", err.Error(), helper.EmptyObj{}))
		return
	}

	res, err := c.service.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, helper.BuildErrorResponse(
			"Satuan produk tidak ditemukan", "NOT_FOUND", "product_unit", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil satuan produk", res))
}

func (c *productUnitController) FindWithPagination(ctx *gin.Context) {
	businessIDStr := ctx.Query("business_id")
	if businessIDStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter business_id wajib diisi", "BAD_REQUEST", "business_id", "Query parameter business_id kosong", helper.EmptyObj{}))
		return
	}

	businessID, err := strconv.Atoi(businessIDStr)
	if err != nil || businessID <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter business_id tidak valid", "BAD_REQUEST", "business_id", err.Error(), helper.EmptyObj{}))
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	sortBy := ctx.DefaultQuery("sortBy", "id")
	orderBy := ctx.DefaultQuery("orderBy", "asc")
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

	productUnits, total, err := c.service.FindWithPagination(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data satuan produk", "INTERNAL_ERROR", "product_unit", err.Error(), helper.EmptyObj{}))
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
		"Berhasil mengambil data satuan produk",
		productUnits,
		paginationMeta,
	))
}
