package controller

import (
	"net/http"

	"run-sync/data/request"
	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RunActivityController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindByUserId(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetUserStats(ctx *gin.Context)
}

type runActivityController struct {
	service service.RunActivityService
}

func NewRunActivityController(s service.RunActivityService) RunActivityController {
	return &runActivityController{service: s}
}

func (c *runActivityController) Create(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	var req request.CreateRunActivityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Create(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal membuat aktivitas lari", "CREATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Aktivitas lari berhasil dibuat", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *runActivityController) Update(ctx *gin.Context) {
	activityId, _ := uuid.Parse(ctx.Param("id"))
	var req request.UpdateRunActivityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Update(activityId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengubah aktivitas lari", "UPDATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Aktivitas lari berhasil diubah", result)
	ctx.JSON(http.StatusOK, response)
}

func (c *runActivityController) FindById(ctx *gin.Context) {
	activityId, _ := uuid.Parse(ctx.Param("id"))
	activity, err := c.service.FindById(activityId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data aktivitas lari", activity)
	ctx.JSON(http.StatusOK, response)
}

func (c *runActivityController) FindByUserId(ctx *gin.Context) {
	userId, _ := uuid.Parse(ctx.Param("userId"))
	activities, err := c.service.FindByUserId(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil aktivitas lari pengguna", activities)
	ctx.JSON(http.StatusOK, response)
}

func (c *runActivityController) FindAll(ctx *gin.Context) {
	activities, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil semua data aktivitas lari", activities)
	ctx.JSON(http.StatusOK, response)
}

func (c *runActivityController) Delete(ctx *gin.Context) {
	activityId, _ := uuid.Parse(ctx.Param("id"))
	err := c.service.Delete(activityId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus aktivitas lari", "DELETE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Aktivitas lari berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
}

func (c *runActivityController) GetUserStats(ctx *gin.Context) {
	userId, _ := uuid.Parse(ctx.Param("userId"))
	stats, err := c.service.GetUserStats(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil statistik pengguna", stats)
	ctx.JSON(http.StatusOK, response)
}
