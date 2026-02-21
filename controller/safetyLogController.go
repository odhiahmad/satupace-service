package controller

import (
	"net/http"

	"run-sync/data/request"
	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SafetyLogController interface {
	Create(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindByUserId(ctx *gin.Context)
	FindByMatchId(ctx *gin.Context)
	FindByStatus(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type safetyLogController struct {
	service service.SafetyLogService
}

func NewSafetyLogController(s service.SafetyLogService) SafetyLogController {
	return &safetyLogController{service: s}
}

func (c *safetyLogController) Create(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	var req request.CreateSafetyLogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Create(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal membuat laporan keamanan", "CREATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Laporan keamanan berhasil dibuat", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *safetyLogController) FindById(ctx *gin.Context) {
	logId, _ := uuid.Parse(ctx.Param("id"))
	log, err := c.service.FindById(logId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data laporan keamanan", log)
	ctx.JSON(http.StatusOK, response)
}

func (c *safetyLogController) FindByUserId(ctx *gin.Context) {
	userId, _ := uuid.Parse(ctx.Param("userId"))
	logs, err := c.service.FindByUserId(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil laporan keamanan pengguna", logs)
	ctx.JSON(http.StatusOK, response)
}

func (c *safetyLogController) FindByMatchId(ctx *gin.Context) {
	matchId, _ := uuid.Parse(ctx.Param("matchId"))
	logs, err := c.service.FindByMatchId(matchId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil laporan keamanan match", logs)
	ctx.JSON(http.StatusOK, response)
}

func (c *safetyLogController) FindByStatus(ctx *gin.Context) {
	status := ctx.Query("status")
	logs, err := c.service.FindByStatus(status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil laporan keamanan berdasarkan status", logs)
	ctx.JSON(http.StatusOK, response)
}

func (c *safetyLogController) Delete(ctx *gin.Context) {
	logId, _ := uuid.Parse(ctx.Param("id"))
	err := c.service.Delete(logId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus laporan keamanan", "DELETE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Laporan keamanan berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
}
