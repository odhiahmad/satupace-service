package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odhiahmad/kasirku-service/data/request"
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
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", err.Error(), nil))
		return
	}

	res, err := c.businessBranchService.Create(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal membuat cabang bisnis", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil membuat cabang bisnis", res))
}

func (c *businessBranchController) Update(ctx *gin.Context) {
	var input request.BusinessBranchUpdate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", err.Error(), nil))
		return
	}

	res, err := c.businessBranchService.Update(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal mengubah cabang bisnis", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengubah cabang bisnis", res))
}

func (c *businessBranchController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", err.Error(), nil))
		return
	}

	err = c.businessBranchService.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal menghapus cabang bisnis", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil menghapus cabang bisnis", nil))
}

func (c *businessBranchController) FindById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("ID tidak valid", err.Error(), nil))
		return
	}

	res, err := c.businessBranchService.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, helper.BuildErrorResponse("Cabang bisnis tidak ditemukan", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil cabang bisnis", res))
}

// ✅ Endpoint: GET /businessBranchs/paginate?business_id=1&page=1&limit=10&search=lebaran
func (c *businessBranchController) FindWithPagination(ctx *gin.Context) {
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

	businessBranches, total, err := c.businessBranchService.FindWithPagination(businessID, pagination)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil data cabang bisnis", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	responseData := gin.H{
		"total":   total,
		"limit":   limit,
		"results": businessBranches,
	}

	res := helper.BuildResponse(true, "Berhasil mengambil data cabang bisnis", responseData)
	ctx.JSON(http.StatusOK, res)
}
