package controller

import (
	"net/http"
	"strconv"

	"run-sync/data/request"
	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SmartWatchController interface {
	ConnectDevice(ctx *gin.Context)
	DisconnectDevice(ctx *gin.Context)
	GetDevices(ctx *gin.Context)
	SyncActivity(ctx *gin.Context)
	BatchSync(ctx *gin.Context)
	GetSyncHistory(ctx *gin.Context)
	GetDeviceStats(ctx *gin.Context)
}

type smartWatchController struct {
	service service.SmartWatchService
}

func NewSmartWatchController(s service.SmartWatchService) SmartWatchController {
	return &smartWatchController{service: s}
}

func (c *smartWatchController) ConnectDevice(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	var req request.ConnectDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.ConnectDevice(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menghubungkan perangkat", "CONNECT_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Perangkat berhasil dihubungkan", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *smartWatchController) DisconnectDevice(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	deviceId, err := uuid.Parse(ctx.Param("device_id"))
	if err != nil {
		res := helper.BuildErrorResponse("ID perangkat tidak valid", "INVALID_ID", "param", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.service.DisconnectDevice(userId, deviceId); err != nil {
		res := helper.BuildErrorResponse("Gagal memutuskan perangkat", "DISCONNECT_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Perangkat berhasil diputuskan", nil)
	ctx.JSON(http.StatusOK, response)
}

func (c *smartWatchController) GetDevices(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	devices, err := c.service.GetDevices(userId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil data perangkat", "FETCH_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data perangkat", devices)
	ctx.JSON(http.StatusOK, response)
}

func (c *smartWatchController) SyncActivity(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	var req request.SyncActivityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.SyncActivity(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menyinkronkan aktivitas", "SYNC_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Aktivitas berhasil disinkronkan", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *smartWatchController) BatchSync(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	var req request.BatchSyncRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.BatchSync(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal batch sync", "BATCH_SYNC_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Batch sync selesai", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *smartWatchController) GetSyncHistory(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	limit := 20
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	history, err := c.service.GetSyncHistory(userId, limit)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil riwayat sync", "FETCH_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil riwayat sync", history)
	ctx.JSON(http.StatusOK, response)
}

func (c *smartWatchController) GetDeviceStats(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	deviceId, err := uuid.Parse(ctx.Param("device_id"))
	if err != nil {
		res := helper.BuildErrorResponse("ID perangkat tidak valid", "INVALID_ID", "param", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	stats, err := c.service.GetDeviceStats(userId, deviceId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil statistik perangkat", "STATS_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil statistik perangkat", stats)
	ctx.JSON(http.StatusOK, response)
}
