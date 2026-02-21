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
	GetCandidates(ctx *gin.Context)
	SendMatchRequest(ctx *gin.Context)
	AcceptMatch(ctx *gin.Context)
	RejectMatch(ctx *gin.Context)
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

func (c *directMatchController) GetCandidates(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	candidates, err := c.service.GetCandidates(userId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengambil kandidat match", "CANDIDATES_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil kandidat match", candidates)
	ctx.JSON(http.StatusOK, response)
}

func (c *directMatchController) SendMatchRequest(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	var req request.CreateDirectMatchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.SendMatchRequest(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengirim match request", "MATCH_REQUEST_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Match request berhasil dikirim", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *directMatchController) AcceptMatch(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	matchId, _ := uuid.Parse(ctx.Param("id"))

	result, err := c.service.AcceptMatch(matchId, userId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menerima match", "ACCEPT_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Match berhasil diterima", result)
	ctx.JSON(http.StatusOK, response)
}

func (c *directMatchController) RejectMatch(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	matchId, _ := uuid.Parse(ctx.Param("id"))

	result, err := c.service.RejectMatch(matchId, userId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menolak match", "REJECT_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Match berhasil ditolak", result)
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
