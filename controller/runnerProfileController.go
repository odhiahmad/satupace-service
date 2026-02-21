package controller

import (
	"net/http"

	"run-sync/data/request"
	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RunnerProfileController interface {
	CreateOrUpdate(ctx *gin.Context)
	Update(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindByUserId(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type runnerProfileController struct {
	service service.RunnerProfileService
}

func NewRunnerProfileController(s service.RunnerProfileService) RunnerProfileController {
	return &runnerProfileController{service: s}
}

func (c *runnerProfileController) CreateOrUpdate(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	var req request.CreateRunnerProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.CreateOrUpdate(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal membuat/memperbarui profil runner", "CREATE_UPDATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Profil runner berhasil disimpan", result)
	ctx.JSON(http.StatusOK, response)
}

func (c *runnerProfileController) Update(ctx *gin.Context) {
	profileId, _ := uuid.Parse(ctx.Param("id"))
	var req request.UpdateRunnerProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Update(profileId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengubah profil runner", "UPDATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Profil runner berhasil diubah", result)
	ctx.JSON(http.StatusOK, response)
}

func (c *runnerProfileController) FindById(ctx *gin.Context) {
	profileId, _ := uuid.Parse(ctx.Param("id"))
	profile, err := c.service.FindById(profileId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data profil runner", profile)
	ctx.JSON(http.StatusOK, response)
}

func (c *runnerProfileController) FindByUserId(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	profile, err := c.service.FindByUserId(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data profil runner", profile)
	ctx.JSON(http.StatusOK, response)
}

func (c *runnerProfileController) FindAll(ctx *gin.Context) {
	profiles, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil semua data profil runner", profiles)
	ctx.JSON(http.StatusOK, response)
}

func (c *runnerProfileController) Delete(ctx *gin.Context) {
	profileId, _ := uuid.Parse(ctx.Param("id"))
	err := c.service.Delete(profileId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus profil runner", "DELETE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Profil runner berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
}
