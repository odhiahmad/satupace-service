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

type BusinessController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindWithPagination(ctx *gin.Context)
}

type businessController struct {
	businessService service.BusinessService
	jwtService      service.JWTService
}

func NewBusinessController(businessService service.BusinessService, jwtService service.JWTService) BusinessController {
	return &businessController{businessService: businessService, jwtService: jwtService}
}

func (c *businessController) Create(ctx *gin.Context) {
	var input request.BusinessCreate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "BAD_REQUEST", "body", err.Error(), nil))
		return
	}

	res, err := c.businessService.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal membuat bisnis", "CREATE_FAILED", "service", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil membuat bisnis", res))
}

func (c *businessController) Update(ctx *gin.Context) {
	var input request.BusinessUpdate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "BAD_REQUEST", "body", err.Error(), nil))
		return
	}

	res, err := c.businessService.Update(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengubah bisnis", "UPDATE_FAILED", "service", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengubah bisnis", res))
}

func (c *businessController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", "BAD_REQUEST", "id", err.Error(), nil))
		return
	}

	err = c.businessService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal menghapus bisnis", "DELETE_FAILED", "service", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil menghapus bisnis", nil))
}

func (c *businessController) FindById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", "BAD_REQUEST", "id", err.Error(), nil))
		return
	}

	res, err := c.businessService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, helper.BuildErrorResponse("Bisnis tidak ditemukan", "NOT_FOUND", "id", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil bisnis", res))
}

func (c *businessController) FindWithPagination(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	sortBy := ctx.DefaultQuery("sortBy", "id")
	orderBy := ctx.DefaultQuery("orderBy", "asc")
	search := ctx.DefaultQuery("search", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter page tidak valid", "INVALID_PARAM", "page", err.Error(), nil))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter limit tidak valid", "INVALID_PARAM", "limit", err.Error(), nil))
		return
	}

	pagination := request.Pagination{
		Page:    page,
		Limit:   limit,
		SortBy:  sortBy,
		OrderBy: orderBy,
		Search:  search,
	}

	businesses, total, err := c.businessService.FindWithPagination(pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengambil data bisnis", "FETCH_FAILED", "service", err.Error(), nil))
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
		"Berhasil mengambil data bisnis",
		businesses,
		paginationMeta,
	))
}
