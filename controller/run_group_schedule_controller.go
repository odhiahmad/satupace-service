package controller

import (
	"net/http"

	"run-sync/data/request"
	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RunGroupScheduleController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	FindByGroupId(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type runGroupScheduleController struct {
	service service.RunGroupScheduleService
}

func NewRunGroupScheduleController(s service.RunGroupScheduleService) RunGroupScheduleController {
	return &runGroupScheduleController{service: s}
}

// Create - POST /runs/groups/:id/schedules
func (c *runGroupScheduleController) Create(ctx *gin.Context) {
	groupId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		res := helper.BuildErrorResponse("ID grup tidak valid", "INVALID_ID", "path", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var req request.CreateRunGroupScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Create(groupId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal membuat jadwal", "CREATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Jadwal berhasil dibuat", result)
	ctx.JSON(http.StatusCreated, response)
}

// Update - PUT /runs/groups/schedules/:scheduleId
func (c *runGroupScheduleController) Update(ctx *gin.Context) {
	scheduleId, err := uuid.Parse(ctx.Param("scheduleId"))
	if err != nil {
		res := helper.BuildErrorResponse("ID jadwal tidak valid", "INVALID_ID", "path", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var req request.UpdateRunGroupScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Update(scheduleId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengubah jadwal", "UPDATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Jadwal berhasil diubah", result)
	ctx.JSON(http.StatusOK, response)
}

// FindByGroupId - GET /runs/groups/:id/schedules
func (c *runGroupScheduleController) FindByGroupId(ctx *gin.Context) {
	groupId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		res := helper.BuildErrorResponse("ID grup tidak valid", "INVALID_ID", "path", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	schedules, err := c.service.FindByGroupId(groupId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil jadwal grup", schedules)
	ctx.JSON(http.StatusOK, response)
}

// Delete - DELETE /runs/groups/schedules/:scheduleId
func (c *runGroupScheduleController) Delete(ctx *gin.Context) {
	scheduleId, err := uuid.Parse(ctx.Param("scheduleId"))
	if err != nil {
		res := helper.BuildErrorResponse("ID jadwal tidak valid", "INVALID_ID", "path", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	if err := c.service.Delete(scheduleId); err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus jadwal", "DELETE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Jadwal berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
}
