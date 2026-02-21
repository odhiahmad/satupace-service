package controller

import (
	"net/http"

	"run-sync/data/request"
	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RunGroupMemberController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindByGroupId(ctx *gin.Context)
	FindByUserId(ctx *gin.Context)
	Delete(ctx *gin.Context)
	JoinGroup(ctx *gin.Context)
	LeaveGroup(ctx *gin.Context)
	KickMember(ctx *gin.Context)
}

type runGroupMemberController struct {
	service service.RunGroupMemberService
}

func NewRunGroupMemberController(s service.RunGroupMemberService) RunGroupMemberController {
	return &runGroupMemberController{service: s}
}

func (c *runGroupMemberController) Create(ctx *gin.Context) {
	var req request.CreateRunGroupMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Create(req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menambahkan anggota grup", "CREATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Anggota grup berhasil ditambahkan", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *runGroupMemberController) Update(ctx *gin.Context) {
	memberId, _ := uuid.Parse(ctx.Param("id"))
	var req request.UpdateRunGroupMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Update(memberId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengubah anggota grup", "UPDATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Anggota grup berhasil diubah", result)
	ctx.JSON(http.StatusOK, response)
}

func (c *runGroupMemberController) FindById(ctx *gin.Context) {
	memberId, _ := uuid.Parse(ctx.Param("id"))
	member, err := c.service.FindById(memberId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data anggota grup", member)
	ctx.JSON(http.StatusOK, response)
}

func (c *runGroupMemberController) FindByGroupId(ctx *gin.Context) {
	groupId, _ := uuid.Parse(ctx.Param("groupId"))
	members, err := c.service.FindByGroupId(groupId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil anggota grup", members)
	ctx.JSON(http.StatusOK, response)
}

func (c *runGroupMemberController) FindByUserId(ctx *gin.Context) {
	userId, _ := uuid.Parse(ctx.Param("userId"))
	members, err := c.service.FindByUserId(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil grup pengguna", members)
	ctx.JSON(http.StatusOK, response)
}

func (c *runGroupMemberController) Delete(ctx *gin.Context) {
	memberId, _ := uuid.Parse(ctx.Param("id"))
	err := c.service.Delete(memberId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus anggota grup", "DELETE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Anggota grup berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
}

func (c *runGroupMemberController) JoinGroup(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	groupId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		res := helper.BuildErrorResponse("ID grup tidak valid", "INVALID_REQUEST", "id", "format UUID tidak valid", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.JoinGroup(userId, groupId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal bergabung dengan grup", "JOIN_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Berhasil bergabung dengan grup", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *runGroupMemberController) UpdateRole(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	memberId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		res := helper.BuildErrorResponse("ID member tidak valid", "INVALID_REQUEST", "id", "format UUID tidak valid", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var req request.UpdateMemberRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.UpdateRole(userId, memberId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengubah role anggota", "UPDATE_ROLE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Role anggota berhasil diubah", result)
	ctx.JSON(http.StatusOK, response)
}

func (c *runGroupMemberController) LeaveGroup(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	groupId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		res := helper.BuildErrorResponse("ID grup tidak valid", "INVALID_REQUEST", "id", "format UUID tidak valid", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err = c.service.LeaveGroup(userId, groupId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal keluar dari grup", "LEAVE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Berhasil keluar dari grup", nil)
	ctx.JSON(http.StatusOK, response)
}

func (c *runGroupMemberController) KickMember(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	memberId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		res := helper.BuildErrorResponse("ID member tidak valid", "INVALID_REQUEST", "id", "format UUID tidak valid", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err = c.service.KickMember(userId, memberId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengeluarkan anggota", "KICK_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Anggota berhasil dikeluarkan dari grup", nil)
	ctx.JSON(http.StatusOK, response)
}
