package controller

import (
	"net/http"

	"run-sync/data/request"
	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DirectMatchController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindUserMatches(ctx *gin.Context)
	FindMatchesByStatus(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type directMatchController struct {
	service service.DirectMatchService
}

func NewDirectMatchController(s service.DirectMatchService) DirectMatchController {
	return &directMatchController{service: s}
}

func (c *directMatchController) Create(ctx *gin.Context) {
	user1Id := ctx.MustGet("user_id").(uuid.UUID)
	var req request.CreateDirectMatchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Create(user1Id, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal membuat match", "CREATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Match berhasil dibuat", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *directMatchController) Update(ctx *gin.Context) {
	matchId, _ := uuid.Parse(ctx.Param("id"))
	var req request.UpdateDirectMatchStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Update(matchId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengubah match", "UPDATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Match berhasil diubah", result)
	ctx.JSON(http.StatusOK, response)
}

func (c *directMatchController) FindById(ctx *gin.Context) {
	matchId, _ := uuid.Parse(ctx.Param("id"))
	match, err := c.service.FindById(matchId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data match", match)
	ctx.JSON(http.StatusOK, response)
}

func (c *directMatchController) FindUserMatches(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	matches, err := c.service.FindUserMatches(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil match pengguna", matches)
	ctx.JSON(http.StatusOK, response)
}

func (c *directMatchController) FindMatchesByStatus(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	status := ctx.Query("status")
	matches, err := c.service.FindMatchesByStatus(userId, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil match berdasarkan status", matches)
	ctx.JSON(http.StatusOK, response)
}

func (c *directMatchController) Delete(ctx *gin.Context) {
	matchId, _ := uuid.Parse(ctx.Param("id"))
	err := c.service.Delete(matchId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus match", "DELETE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Match berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
}
