package controller

import (
	"net/http"

	"run-sync/data/request"
	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserPhotoController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindByUserId(ctx *gin.Context)
	FindMyPhotos(ctx *gin.Context)
	FindPrimaryPhoto(ctx *gin.Context)
	Delete(ctx *gin.Context)
	VerifyFace(ctx *gin.Context)
}

type userPhotoController struct {
	service service.UserPhotoService
}

func NewUserPhotoController(s service.UserPhotoService) UserPhotoController {
	return &userPhotoController{service: s}
}

func (c *userPhotoController) Create(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	var req request.UploadUserPhotoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Create(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengunggah foto", "CREATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Foto berhasil diunggah", result)
	ctx.JSON(http.StatusCreated, response)
}

func (c *userPhotoController) Update(ctx *gin.Context) {
	photoId, _ := uuid.Parse(ctx.Param("id"))
	var req request.UpdateUserPhotoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.Update(photoId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengubah foto", "UPDATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Foto berhasil diubah", result)
	ctx.JSON(http.StatusOK, response)
}

func (c *userPhotoController) FindById(ctx *gin.Context) {
	photoId, _ := uuid.Parse(ctx.Param("id"))
	photo, err := c.service.FindById(photoId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil data foto", photo)
	ctx.JSON(http.StatusOK, response)
}

func (c *userPhotoController) FindByUserId(ctx *gin.Context) {
	userId, _ := uuid.Parse(ctx.Param("userId"))
	photos, err := c.service.FindByUserId(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil foto pengguna", photos)
	ctx.JSON(http.StatusOK, response)
}

func (c *userPhotoController) FindMyPhotos(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	photos, err := c.service.FindByUserId(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil foto saya", photos)
	ctx.JSON(http.StatusOK, response)
}

func (c *userPhotoController) FindPrimaryPhoto(ctx *gin.Context) {
	userId, _ := uuid.Parse(ctx.Param("userId"))
	photo, err := c.service.FindPrimaryPhoto(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := helper.BuildResponse(true, "Berhasil mengambil foto utama", photo)
	ctx.JSON(http.StatusOK, response)
}

func (c *userPhotoController) Delete(ctx *gin.Context) {
	photoId, _ := uuid.Parse(ctx.Param("id"))
	err := c.service.Delete(photoId)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal menghapus foto", "DELETE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Foto berhasil dihapus", nil)
	ctx.JSON(http.StatusOK, response)
}

func (c *userPhotoController) VerifyFace(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	var req request.FaceVerifyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.service.VerifyFace(userId, req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal verifikasi wajah", "VERIFY_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := helper.BuildResponse(true, "Verifikasi wajah selesai", result)
	ctx.JSON(http.StatusOK, response)
}
