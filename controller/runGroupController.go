package controller

import (
	"net/http"

	"run-sync/data/request"
	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RunGroupController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	FindByStatus(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindByCreatedBy(ctx *gin.Context)
}

type runGroupController struct {
	service service.RunGroupService
}

func NewRunGroupController(s service.RunGroupService) RunGroupController {
	return &runGroupController{service: s}
}

func (c *runGroupController) Create(ctx *gin.Context) {
	createdBy := ctx.MustGet("user_id").(uuid.UUID)
	var req request.CreateRunGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Create(createdBy, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal membuat grup lari", "CREATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Grup lari berhasil dibuat", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *runGroupController) Update(ctx *gin.Context) {
	groupId, _ := uuid.Parse(ctx.Param("id"))
	var req request.UpdateRunGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Update(groupId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengubah grup lari", "UPDATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Grup lari berhasil diubah", result)
	ctx.JSON(http.StatusOK, response)
}

func (c *runGroupController) FindById(ctx *gin.Context) {
	groupId, _ := uuid.Parse(ctx.Param("id"))
	group, err := c.service.FindById(groupId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data grup lari", group)
	ctx.JSON(http.StatusOK, response)
}

func (c *runGroupController) FindAll(ctx *gin.Context) {
	groups, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil semua data grup lari", groups)
	ctx.JSON(http.StatusOK, response)
}

func (c *runGroupController) FindByStatus(ctx *gin.Context) {
	status := ctx.Query("status")
	groups, err := c.service.FindByStatus(status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data grup lari", groups)
	ctx.JSON(http.StatusOK, response)
}

func (c *runGroupController) Delete(ctx *gin.Context) {
	groupId, _ := uuid.Parse(ctx.Param("id"))
	err := c.service.Delete(groupId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus grup lari", "DELETE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Grup lari berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
}

func (c *runGroupController) FindByCreatedBy(ctx *gin.Context) {
	userId, _ := uuid.Parse(ctx.Param("userId"))
	groups, err := c.service.FindByCreatedBy(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data grup lari", groups)
	ctx.JSON(http.StatusOK, response)
}
