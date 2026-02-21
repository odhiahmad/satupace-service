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

type StravaController interface {
	GetAuthURL(ctx *gin.Context)
	Callback(ctx *gin.Context)
	Disconnect(ctx *gin.Context)
	GetConnection(ctx *gin.Context)
	SyncActivities(ctx *gin.Context)
	GetSyncHistory(ctx *gin.Context)
	GetStats(ctx *gin.Context)
}

type stravaController struct {
	service service.StravaService
}

func NewStravaController(s service.StravaService) StravaController {
	return &stravaController{service: s}
}

// GetAuthURL returns the Strava OAuth authorization URL
func (c *stravaController) GetAuthURL(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	result, err := c.service.GetAuthURL(userId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal membuat URL autentikasi Strava", "STRAVA_AUTH_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "URL autentikasi Strava berhasil dibuat", result)
	ctx.JSON(http.StatusOK, response)
}

// Callback handles the OAuth callback from Strava
func (c *stravaController) Callback(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	var req request.StravaCallbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.HandleCallback(userId, req.Code)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menghubungkan Strava", "STRAVA_CONNECT_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Strava berhasil dihubungkan", result)
	ctx.JSON(http.StatusOK, response)
}

// Disconnect removes the Strava connection
func (c *stravaController) Disconnect(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	if err := c.service.Disconnect(userId); err != nil {
		res := helper.BuildErrorResponse("Gagal memutuskan Strava", "STRAVA_DISCONNECT_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Strava berhasil diputuskan", nil)
	ctx.JSON(http.StatusOK, response)
}

// GetConnection returns the current Strava connection status
func (c *stravaController) GetConnection(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	result, err := c.service.GetConnection(userId)
	if err != nil {
		res := helper.BuildErrorResponse("Belum terhubung dengan Strava", "NOT_CONNECTED", "body", err.Error(), nil)
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data koneksi Strava", result)
	ctx.JSON(http.StatusOK, response)
}

// SyncActivities fetches and syncs running activities from Strava
func (c *stravaController) SyncActivities(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	result, err := c.service.SyncActivities(userId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal sinkronisasi aktivitas Strava", "STRAVA_SYNC_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Sinkronisasi aktivitas Strava selesai", result)
	ctx.JSON(http.StatusOK, response)
}

// GetSyncHistory returns synced Strava activities
func (c *stravaController) GetSyncHistory(ctx *gin.Context) {
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

	response := helper.BuildResponse(true, "Berhasil mengambil riwayat sync Strava", history)
	ctx.JSON(http.StatusOK, response)
}

// GetStats returns Strava sync statistics
func (c *stravaController) GetStats(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	stats, err := c.service.GetStats(userId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil statistik Strava", "STATS_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil statistik Strava", stats)
	ctx.JSON(http.StatusOK, response)
}
