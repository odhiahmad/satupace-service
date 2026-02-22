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

type NotificationController interface {
	GetMyNotifications(ctx *gin.Context)
	MarkAsRead(ctx *gin.Context)
	MarkAllAsRead(ctx *gin.Context)
	RegisterDeviceToken(ctx *gin.Context)
	RemoveDeviceToken(ctx *gin.Context)
}

type notificationController struct {
	service service.NotificationService
}

func NewNotificationController(s service.NotificationService) NotificationController {
	return &notificationController{service: s}
}

// GET /notifications
func (c *notificationController) GetMyNotifications(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	result, err := c.service.GetByUser(userId, page, limit)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil notifikasi", "NOTIF_FETCH_FAILED", "server", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Berhasil mengambil notifikasi", result))
}

// PATCH /notifications/read
func (c *notificationController) MarkAsRead(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	var req request.MarkNotificationReadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.service.MarkAsRead(userId, req.Ids); err != nil {
		res := helper.BuildErrorResponse("Gagal menandai notifikasi", "MARK_READ_FAILED", "server", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Notifikasi berhasil ditandai sudah dibaca", nil))
}

// PATCH /notifications/read-all
func (c *notificationController) MarkAllAsRead(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	if err := c.service.MarkAllAsRead(userId); err != nil {
		res := helper.BuildErrorResponse("Gagal menandai semua notifikasi", "MARK_ALL_READ_FAILED", "server", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Semua notifikasi berhasil ditandai sudah dibaca", nil))
}

// POST /notifications/device-token
func (c *notificationController) RegisterDeviceToken(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	var req request.RegisterDeviceTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.service.RegisterDeviceToken(userId, req.FCMToken, req.Platform); err != nil {
		res := helper.BuildErrorResponse("Gagal menyimpan device token", "TOKEN_SAVE_FAILED", "server", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Device token berhasil disimpan", nil))
}

// DELETE /notifications/device-token
func (c *notificationController) RemoveDeviceToken(ctx *gin.Context) {
	var req request.RemoveDeviceTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.service.RemoveDeviceToken(req.FCMToken); err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus device token", "TOKEN_DELETE_FAILED", "server", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusOK, helper.BuildResponse(true, "Device token berhasil dihapus", nil))
}
