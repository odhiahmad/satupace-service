package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type ProductController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	Delete(ctx *gin.Context)
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

// Create Product
func (c *productController) Create(ctx *gin.Context) {
	var req request.ProductCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Gagal bind data", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.productService.Create(req); err != nil {
		res := helper.BuildErrorResponse("Gagal membuat produk", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Produk berhasil dibuat", helper.EmptyObj{})
	ctx.JSON(http.StatusCreated, res)
}

// Update Product
func (c *productController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse("ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var req request.ProductUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Gagal bind data", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.productService.Update(id, req); err != nil {
		res := helper.BuildErrorResponse("Gagal mengupdate produk", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Produk berhasil diperbarui", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

// Find By ID
func (c *productController) FindById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse("ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	product, err := c.productService.FindById(id)
	if err != nil {
		res := helper.BuildErrorResponse("Produk tidak ditemukan", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	res := helper.BuildResponse(true, "Produk berhasil ditemukan", product)
	ctx.JSON(http.StatusOK, res)
}

// Find All
func (c *productController) FindAll(ctx *gin.Context) {
	products, err := c.productService.FindAll()
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil data produk", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := helper.BuildResponse(true, "Semua produk berhasil diambil", products)
	ctx.JSON(http.StatusOK, res)
}

// Delete Product
func (c *productController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res := helper.BuildErrorResponse("ID tidak valid", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.productService.Delete(id); err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus produk", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := helper.BuildResponse(true, "Produk berhasil dihapus", helper.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}
