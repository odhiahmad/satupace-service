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

type TransactionController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	AddOrUpdateItem(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindWithPagination(ctx *gin.Context)
}

type transactionController struct {
	transactionService service.TransactionService
	jwtService         service.JWTService
}

func NewTransactionController(transactionService service.TransactionService, jwtService service.JWTService) TransactionController {
	return &transactionController{
		transactionService: transactionService,
		jwtService:         jwtService,
	}
}

func (c *transactionController) Create(ctx *gin.Context) {
	var input request.TransactionCreateRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "BAD_REQUEST", "body", err.Error(), nil))
		return
	}

	transaction, err := c.transactionService.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal membuat transaksi", "INTERNAL_ERROR", "transaction", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, helper.BuildResponse(true, "Berhasil membuat transaksi", transaction))
}

func (c *transactionController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", "BAD_REQUEST", "id", err.Error(), nil))
		return
	}

	var input request.TransactionUpdateRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "BAD_REQUEST", "body", err.Error(), nil))
		return
	}

	transaction, err := c.transactionService.Update(id, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengubah transaksi", "INTERNAL_ERROR", "transaction", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengubah transaksi", transaction))
}

func (c *transactionController) AddOrUpdateItem(ctx *gin.Context) {
	transactionId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || transactionId <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID transaksi tidak valid", "BAD_REQUEST", "id", err.Error(), nil))
		return
	}

	var input request.TransactionItemCreate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input item tidak valid", "BAD_REQUEST", "body", err.Error(), nil))
		return
	}

	transaction, err := c.transactionService.AddOrUpdateItem(transactionId, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal menambahkan item", "INTERNAL_ERROR", "transaction_item", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Item berhasil ditambahkan/diupdate", transaction))
}

func (c *transactionController) FindById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", "BAD_REQUEST", "id", err.Error(), nil))
		return
	}

	res, err := c.transactionService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, helper.BuildErrorResponse("Transaksi tidak ditemukan", "NOT_FOUND", "transaction", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil transaksi", res))
}

func (c *transactionController) FindWithPagination(ctx *gin.Context) {
	businessIDStr := ctx.Query("business_id")
	if businessIDStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter business_id wajib diisi", "BAD_REQUEST", "business_id", "business_id query param is required", nil))
		return
	}

	businessID, err := strconv.Atoi(businessIDStr)
	if err != nil || businessID <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter business_id tidak valid", "BAD_REQUEST", "business_id", err.Error(), nil))
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	sortBy := ctx.DefaultQuery("sortBy", "id")
	orderBy := ctx.DefaultQuery("orderBy", "desc")
	search := ctx.DefaultQuery("search", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter page tidak valid", "BAD_REQUEST", "page", err.Error(), nil))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter limit tidak valid", "BAD_REQUEST", "limit", err.Error(), nil))
		return
	}

	pagination := request.Pagination{
		Page:    page,
		Limit:   limit,
		SortBy:  sortBy,
		OrderBy: orderBy,
		Search:  search,
	}

	transactions, total, err := c.transactionService.FindWithPagination(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengambil data transaksi", "INTERNAL_ERROR", "transaction", err.Error(), nil))
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	paginationMeta := response.PaginatedResponse{
		Page:      page,
		Limit:     limit,
		Total:     total,
		OrderBy:   sortBy,
		SortOrder: orderBy,
	}

	response := gin.H{
		"results":    transactions,
		"total":      total,
		"limit":      limit,
		"page":       page,
		"totalPages": totalPages,
	}

	ctx.JSON(http.StatusOK, helper.BuildResponsePagination(
		true,
		"Berhasil mengambil data transaksi",
		response,
		paginationMeta,
	))
}
