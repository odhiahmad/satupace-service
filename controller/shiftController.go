package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/odhiahmad/kasirku-service/data/request"
	"github.com/odhiahmad/kasirku-service/helper"
	"github.com/odhiahmad/kasirku-service/service"
)

type ShiftController interface {
	OpenShift(ctx *gin.Context)
	CloseShift(ctx *gin.Context)
	GetActiveShift(ctx *gin.Context)
}

type shiftController struct {
	shiftService service.ShiftService
	jwtService   service.JWTService
}

func NewShiftController(shiftService service.ShiftService, jwtService service.JWTService) ShiftController {
	return &shiftController{shiftService: shiftService, jwtService: jwtService}
}

func (c *shiftController) OpenShift(ctx *gin.Context) {
	var req request.OpenShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "bad_request", "body", err.Error(), nil))
		return
	}

	shift, err := c.shiftService.OpenShift(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Gagal membuka shift", "bad_request", "shift", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, helper.BuildResponse(true, "Shift berhasil dibuka", shift))
}

func (c *shiftController) CloseShift(ctx *gin.Context) {
	var req request.CloseShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Input tidak valid", "bad_request", "body", err.Error(), nil))
		return
	}

	shift, err := c.shiftService.CloseShift(req)
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
