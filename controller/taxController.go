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
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "bad_request", "body", err.Error(), nil))
		return
	}

	res, err := c.taxService.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal membuat tax", "internal_error", "tax", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil membuat tax", res))
}

func (c *taxController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter id wajib diisi", "missing_parameter", "id", "parameter id kosong", nil))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter id tidak valid", "invalid_parameter", "id", err.Error(), nil))
		return
	}

	var input request.TaxUpdate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "bad_request", "body", err.Error(), nil))
		return
	}

	res, err := c.taxService.Update(id, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengubah tax", "internal_error", "tax", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengubah tax", res))
}

func (c *taxController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", "invalid_parameter", "id", err.Error(), nil))
		return
	}

	err = c.taxService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal menghapus tax", "internal_error", "tax", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil menghapus tax", nil))
}

func (c *taxController) FindById(ctx *gin.Context) {
	taxIdParam := ctx.Param("id")
	taxId, err := strconv.Atoi(taxIdParam)
	if err != nil {
		response := helper.BuildErrorResponse(
			"Parameter id tidak valid",
			"invalid_parameter",
			"id",
			err.Error(),
			helper.EmptyObj{},
		)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// hanya satu return value
	taxResponse := c.taxService.FindById(taxId)

	response := helper.BuildResponse(true, "Berhasil mengambil data tax", taxResponse)
	ctx.JSON(http.StatusOK, response)
}

func (c *taxController) FindWithPagination(ctx *gin.Context) {
	businessIDStr := ctx.Query("business_id")
	if businessIDStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter business_id wajib diisi", "missing_parameter", "business_id", "parameter business_id kosong", nil))
		return
	}

	businessID, err := strconv.Atoi(businessIDStr)
	if err != nil || businessID <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter business_id tidak valid", "invalid_parameter", "business_id", err.Error(), nil))
		return
	}

	limitStr := ctx.DefaultQuery("limit", "10")
	sortBy := ctx.DefaultQuery("sortBy", "created_at")
	orderBy := ctx.DefaultQuery("orderBy", "desc")
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

	taxes, total, err := c.taxService.FindWithPagination(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengambil data tax", "internal_error", "tax", err.Error(), nil))
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
		"Data tax berhasil diambil",
		taxes,
		paginationMeta,
	))
}
