package controller

import (
	"context"
	"net/http"
	"time"

	"run-sync/helper"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type WhatsAppController interface {
	Register(ctx *gin.Context)
	Verify(ctx *gin.Context)
}

type whatsappController struct {
	emailHelper *helper.EmailHelper
	redisClient *redis.Client
}

func NewWhatsAppController(emailHelper *helper.EmailHelper, redisClient *redis.Client) WhatsAppController {
	return &whatsappController{emailHelper: emailHelper, redisClient: redisClient}
}

type waRegisterRequest struct {
	Phone string `json:"phone" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type waVerifyRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

func (w *whatsappController) Register(c *gin.Context) {
	var req waRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil))
		return
	}

	// Generate short OTP and save mapping phone->email with OTP hash in Redis
	otp := helper.GenerateOTPCode(6)
	data := map[string]string{"email": req.Email, "otp_hash": helper.HashOTP(otp)}
	key := "whatsapp:register:" + req.Phone

	// store for 24 hours
	if err := helper.SetJSONToRedis(context.Background(), w.redisClient, key, data, 24*time.Hour); err != nil {
		c.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal menyimpan data pendaftaran WhatsApp", "SAVE_FAILED", "body", err.Error(), nil))
		return
	}

	// Send OTP message via WhatsApp if client available
	msg := "Kode verifikasi WhatsApp Anda: " + otp
	if err := helper.SendOTPViaWhatsApp(req.Phone, msg); err != nil {
		// still return success but warn user that message could not be sent
		c.JSON(http.StatusOK, helper.BuildResponse(false, "Pendaftaran tersimpan, namun pengiriman WhatsApp gagal", gin.H{"warning": err.Error()}))
		return
	}

	c.JSON(http.StatusCreated, helper.BuildResponse(true, "Kode verifikasi terkirim ke WhatsApp", nil))
}

func (w *whatsappController) Verify(c *gin.Context) {
	var req waVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil))
		return
	}

	key := "whatsapp:register:" + req.Phone
	var stored map[string]string
	if err := helper.GetJSONFromRedis(context.Background(), w.redisClient, key, &stored); err != nil {
		c.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Data pendaftaran tidak ditemukan atau sudah kadaluarsa", "NOT_FOUND", "body", err.Error(), nil))
		return
	}

	if stored == nil {
		c.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Data pendaftaran tidak ditemukan", "NOT_FOUND", "body", "", nil))
		return
	}

	otpHash, ok := stored["otp_hash"]
	if !ok || otpHash != helper.HashOTP(req.Code) {
		c.JSON(http.StatusBadRequest, helper.BuildErrorResponse("Kode verifikasi tidak valid", "INVALID_CODE", "body", "", nil))
		return
	}

	// mark verified mapping
	email := stored["email"]
	verifiedKey := "whatsapp:verified:" + req.Phone
	if err := w.redisClient.Set(context.Background(), verifiedKey, email, 0).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, helper.BuildErrorResponse("Gagal menyimpan verifikasi", "SAVE_FAILED", "body", err.Error(), nil))
		return
	}

	// remove the registration key
	_ = w.redisClient.Del(context.Background(), key).Err()

	c.JSON(http.StatusOK, helper.BuildResponse(true, "WhatsApp berhasil terverifikasi dan terhubung ke email", gin.H{"email": email}))
}
