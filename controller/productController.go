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

type ProductController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	UpdateImage(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindWithPagination(ctx *gin.Context)
	SetActive(ctx *gin.Context)
	SetAvailable(ctx *gin.Context)
}

type productController struct {
	productService service.ProductService
	jwtService     service.JWTService
}

func NewProductController(productService service.ProductService, jwtService service.JWTService) ProductController {
	return &productController{
		productService: productService,
		jwtService:     jwtService,
	}
}

// CREATE PRODUCT
func (c *productController) Create(ctx *gin.Context) {
	businessId := ctx.MustGet("business_id").(int)
	var input request.ProductRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Gagal bind data product", "INVALID_JSON", "body", err.Error(), helper.EmptyObj{}))
		return
	}

	input.BusinessId = businessId
	res, err := c.productService.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal membuat product", "internal_error", "product", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, helper.BuildResponse(true, "Berhasil membuat product", res))
}

func (c *productController) Update(ctx *gin.Context) {
	businessId := ctx.MustGet("business_id").(int)
	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID produk tidak valid", "INVALID_ID", "id", err.Error(), helper.EmptyObj{}))
		return
	}

	var input request.ProductRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Gagal bind data update", "INVALID_JSON", "body", err.Error(), helper.EmptyObj{}))
		return
	}

	input.BusinessId = businessId
	result, err := c.productService.Update(id, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengubah produk", "UPDATE_ERROR", "product", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Produk berhasil diubah", result))
}

func (c *productController) UpdateImage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse("ID tidak valid", "INVALID_ID", "param", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var req struct {
		Image string `json:"image" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Gambar tidak valid", "INVALID_IMAGE", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	updated, err := c.productService.UpdateImage(id, req.Image)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal update gambar produk", "UPDATE_IMAGE_FAILED", "service", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Gambar Produk berhasil diubah", updated))
}

func (c *productController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID produk tidak valid", "INVALID_ID", "id", err.Error(), helper.EmptyObj{}))
		return
	}

	if err := c.productService.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal menghapus produk", "DELETE_ERROR", "product", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Produk berhasil dihapus", helper.EmptyObj{}))
}

func (c *productController) FindById(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID produk tidak valid", "INVALID_ID", "id", err.Error(), helper.EmptyObj{}))
		return
	}

	product, err := c.productService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, helper.BuildErrorResponse(
			"Produk tidak ditemukan", "NOT_FOUND", "product", err.Error(), helper.EmptyObj{}))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Produk ditemukan", product))
}

func (c *productController) FindWithPagination(ctx *gin.Context) {
	businessID := ctx.MustGet("business_id").(int)
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	orderBy := ctx.DefaultQuery("order_by", "desc")
	search := ctx.DefaultQuery("search", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter page tidak valid", "INVALID_PARAM", "page", err.Error(), helper.EmptyObj{}))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter limit tidak valid", "INVALID_PARAM", "limit", err.Error(), helper.EmptyObj{}))
		return
	}

	pagination := request.Pagination{
		Page:    page,
		Limit:   limit,
		SortBy:  sortBy,
		OrderBy: orderBy,
		Search:  search,
	}

	products, total, err := c.productService.FindWithPagination(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data produk", "FETCH_ERROR", "product", err.Error(), helper.EmptyObj{}))
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
		"Data produk berhasil diambil",
		products,
		paginationMeta,
	))
}

func (c *productController) SetActive(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID tidak valid", "INVALID_ID", "id", err.Error(), helper.EmptyObj{}))
		return
	}

	var body struct {
		IsActive bool `json:"is_active"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Gagal bind body", "INVALID_JSON", "body", err.Error(), helper.EmptyObj{}))
		return
	}

	if err := c.productService.SetActive(id, body.IsActive); err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengubah status aktif produk", "UPDATE_ERROR", "is_active", err.Error(), helper.EmptyObj{}))
		return
	}

	statusMsg := "dinonaktifkan"
	if body.IsActive {
		statusMsg = "diaktifkan"
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Status produk berhasil "+statusMsg, helper.EmptyObj{}))
}

func (c *productController) SetAvailable(ctx *gin.Context) {
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

	if err := c.productService.SetAvailable(id, body.IsAvailable); err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengubah status ketersediaan produk", "UPDATE_ERROR", "is_available", err.Error(), helper.EmptyObj{}))
		return
	}

	statusMsg := "tidak tersedia"
	if body.IsAvailable {
		statusMsg = "tersedia"
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Status produk sekarang ini "+statusMsg, helper.EmptyObj{}))
}
