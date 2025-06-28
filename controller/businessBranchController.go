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

type BusinessBranchController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindWithPagination(ctx *gin.Context)
}

type businessBranchController struct {
	businessBranchService service.BusinessBranchService
	jwtService            service.JWTService
}

func NewBusinessBranchController(businessBranchService service.BusinessBranchService, jwtService service.JWTService) BusinessBranchController {
	return &businessBranchController{
		businessBranchService: businessBranchService,
		jwtService:            jwtService,
	}
}

func (c *businessBranchController) Create(ctx *gin.Context) {
	var input request.BusinessBranchCreate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "BAD_REQUEST", "body", err.Error(), nil))
		return
	}

	res, err := c.businessBranchService.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal membuat cabang bisnis", "CREATE_FAILED", "service", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil membuat cabang bisnis", res))
}

func (c *businessBranchController) Update(ctx *gin.Context) {
	var input request.BusinessBranchUpdate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "BAD_REQUEST", "body", err.Error(), nil))
		return
	}

	res, err := c.businessBranchService.Update(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengubah cabang bisnis", "UPDATE_FAILED", "service", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengubah cabang bisnis", res))
}

func (c *businessBranchController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", "BAD_REQUEST", "id", err.Error(), nil))
		return
	}

	err = c.businessBranchService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal menghapus cabang bisnis", "DELETE_FAILED", "service", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil menghapus cabang bisnis", nil))
}

func (c *businessBranchController) FindById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", "BAD_REQUEST", "id", err.Error(), nil))
		return
	}

	res, err := c.businessBranchService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, helper.BuildErrorResponse("Cabang bisnis tidak ditemukan", "NOT_FOUND", "id", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil cabang bisnis", res))
}

func (c *businessBranchController) FindWithPagination(ctx *gin.Context) {
	businessIDStr := ctx.Query("business_id")
	if businessIDStr == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter business_id wajib diisi", "MISSING_PARAM", "business_id", "parameter business_id tidak ditemukan", nil))
		return
	}

	businessID, err := strconv.Atoi(businessIDStr)
	if err != nil || businessID <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter business_id tidak valid", "INVALID_PARAM", "business_id", err.Error(), nil))
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
			"Parameter page tidak valid", "INVALID_PARAM", "page", err.Error(), nil))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter limit tidak valid", "INVALID_PARAM", "limit", err.Error(), nil))
		return
	}

	pagination := request.Pagination{
		Page:    page,
		Limit:   limit,
		SortBy:  sortBy,
		OrderBy: orderBy,
		Search:  search,
	}

	businessBranches, total, err := c.businessBranchService.FindWithPagination(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data cabang bisnis", "FETCH_FAILED", "service", err.Error(), nil))
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
		"Berhasil mengambil data cabang bisnis",
		businessBranches,
		paginationMeta,
	))
}
