package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/data/response"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type ShiftController interface {
	OpenShift(ctx *gin.Context)
	CloseShift(ctx *gin.Context)
	GetActiveShift(ctx *gin.Context)
	FindWithPaginationCursor(ctx *gin.Context)
}

type shiftController struct {
	shiftService service.ShiftService
	jwtService   service.JWTService
}

func NewShiftController(shiftService service.ShiftService, jwtService service.JWTService) ShiftController {
	return &shiftController{shiftService: shiftService, jwtService: jwtService}
}

func (c *shiftController) OpenShift(ctx *gin.Context) {
	businessIdStr := ctx.MustGet("business_id").(string)
	businessId, err := uuid.Parse(businessIdStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid business_id UUID"})
		return
	}

	cashierIdStr := ctx.MustGet("user_id").(string)
	cashierId, err := uuid.Parse(cashierIdStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id UUID"})
		return
	}

	var req request.OpenShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "bad_request", "body", err.Error(), nil))
		return
	}

	req.BusinessId = businessId
	req.CashierId = cashierId

	shift, err := c.shiftService.OpenShift(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Gagal membuka shift", "bad_request", "shift", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, helper.BuildResponse(true, "Shift berhasil dibuka", shift))
}

func (c *shiftController) CloseShift(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"ID produk tidak valid", "INVALID_ID", "id", err.Error(), helper.EmptyObj{}))
		return
	}

	var req request.CloseShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "bad_request", "body", err.Error(), nil))
		return
	}

	shift, err := c.shiftService.CloseShift(id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Gagal menutup shift", "bad_request", "shift", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Shift berhasil ditutup", shift))
}

func (c *shiftController) GetActiveShift(ctx *gin.Context) {
	terminalId := ctx.Param("terminal_id")
	if terminalId == "" {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Parameter terminal_id wajib diisi", "missing_parameter", "terminal_id", "terminal_id kosong", nil))
		return
	}

	if _, err := uuid.Parse(terminalId); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("terminal_id tidak valid", "invalid_parameter", "terminal_id", err.Error(), nil))
		return
	}

	shift, err := c.shiftService.GetActiveShift(terminalId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, helper.BuildErrorResponse("Tidak ada shift aktif", "not_found", "shift", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Shift aktif ditemukan", shift))
}

func (c *shiftController) FindWithPaginationCursor(ctx *gin.Context) {
	businessIdStr := ctx.MustGet("business_id").(string)
	businessID, err := uuid.Parse(businessIdStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid business_id UUID"})
		return
	}
	limitStr := ctx.DefaultQuery("limit", "10")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	orderBy := ctx.DefaultQuery("order_by", "desc")
	search := ctx.DefaultQuery("search", "")
	cursor := ctx.DefaultQuery("cursor", "")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse(
			"Parameter limit tidak valid", "invalid_parameter", "limit", err.Error(), nil))
		return
	}

	pagination := request.Pagination{
		Cursor:  cursor,
		Limit:   limit,
		SortBy:  sortBy,
		OrderBy: orderBy,
		Search:  search,
	}

	shiftes, nextCursor, hasNext, err := c.shiftService.FindWithPaginationCursor(businessID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helper.BuildErrorResponse(
			"Gagal mengambil data shift", "internal_error", "shift", err.Error(), nil))
		return
	}

	paginationMeta := response.CursorPaginatedResponse{
		Limit:      limit,
		SortBy:     sortBy,
		OrderBy:    orderBy,
		NextCursor: nextCursor,
		HasNext:    hasNext,
	}

	ctx.JSON(http.StatusOK, helper.BuildResponseCursorPagination(
		true,
		"Data shift berhasil diambil",
		shiftes,
		paginationMeta,
	))
}
