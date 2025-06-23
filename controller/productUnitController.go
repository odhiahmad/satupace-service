package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
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
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", err.Error(), nil))
		return
	}

	res, err := c.service.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal membuat diskon", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil membuat diskon", res))
}

func (c *productUnitController) Update(ctx *gin.Context) {
	var input request.ProductUnitUpdate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", err.Error(), nil))
		return
	}

	res, err := c.service.Update(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengubah diskon", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengubah diskon", res))
}

func (c *productUnitController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", err.Error(), nil))
		return
	}

	err = c.service.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal menghapus diskon", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil menghapus diskon", nil))
}

func (c *productUnitController) FindById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", err.Error(), nil))
		return
	}

	res, err := c.service.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, helper.BuildErrorResponse("Diskon tidak ditemukan", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil diskon", res))
}

func (c *productUnitController) FindWithPagination(ctx *gin.Context) {
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

	productUnits, total, err := c.service.FindWithPagination(businessID, pagination)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil data satuan produk", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	response := gin.H{
		"total":      total,
		"limit":      limit,
		"results":    productUnits,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil data satuan produk", response))
}
