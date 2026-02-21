package controller

import (
	"net/http"

	"run-sync/data/request"
	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupChatMessageController interface {
	Create(ctx *gin.Context)
	FindByGroupId(ctx *gin.Context)
	FindBySenderId(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type groupChatMessageController struct {
	service service.GroupChatMessageService
}

func NewGroupChatMessageController(s service.GroupChatMessageService) GroupChatMessageController {
	return &groupChatMessageController{service: s}
}

func (c *groupChatMessageController) Create(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	var req request.SendGroupChatMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Create(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengirim pesan", "CREATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Pesan berhasil dikirim", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *groupChatMessageController) FindByGroupId(ctx *gin.Context) {
	groupId, _ := uuid.Parse(ctx.Param("groupId"))
	messages, err := c.service.FindByGroupId(groupId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil pesan grup", messages)
	ctx.JSON(http.StatusOK, response)
}

func (c *groupChatMessageController) FindBySenderId(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	messages, err := c.service.FindBySenderId(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil pesan yang dikirim", messages)
	ctx.JSON(http.StatusOK, response)
}

func (c *groupChatMessageController) Delete(ctx *gin.Context) {
	messageId, _ := uuid.Parse(ctx.Param("id"))
	err := c.service.Delete(messageId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus pesan", "DELETE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Pesan berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
}
