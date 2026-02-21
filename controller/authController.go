package controller

import (
	"net/http"
	"time"

	"run-sync/data/request"
	"run-sync/helper"
	"run-sync/service"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(ctx *gin.Context)
	VerifyOTP(ctx *gin.Context)
	Login(ctx *gin.Context)
	ResendOTP(ctx *gin.Context)
}

type authController struct {
	userService service.UserService
	jwtService  service.JWTService
	otpHelper   *helper.RedisHelper
	emailHelper *helper.EmailHelper
}

func NewAuthController(
	userService service.UserService,
	jwtService service.JWTService,
	otpHelper *helper.RedisHelper,
	emailHelper *helper.EmailHelper,
) AuthController {
	return &authController{
		userService: userService,
		jwtService:  jwtService,
		otpHelper:   otpHelper,
		emailHelper: emailHelper,
	}
}

// Register - Create user and send OTP for verification
func (c *authController) Register(ctx *gin.Context) {
	var req request.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Create user (isVerified = false, isActive = false)
	result, err := c.userService.Create(req)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal membuat user", "CREATE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Generate and send OTP
	otp := helper.GenerateOTPCode(6)
	hashedOTP := helper.HashOTP(otp)

	// Check rate limit for OTP
	identifier := req.PhoneNumber
	if err := c.otpHelper.AllowRequest(identifier, 5, 15*time.Minute); err != nil {
		res := helper.BuildErrorResponse("Rate limit exceeded", "RATE_LIMIT", "body", err.Error(), nil)
		ctx.JSON(http.StatusTooManyRequests, res)
		return
	}

	// Save OTP to Redis (15 minutes expiry)
	if err := c.otpHelper.SaveOTP("register", identifier, hashedOTP, 15*time.Minute); err != nil {
		res := helper.BuildErrorResponse("Gagal menyimpan OTP", "OTP_SAVE_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// Send OTP via email or WhatsApp
	// For now, we'll return it in response (in production, send via email/SMS)
	response := helper.BuildResponse(true, "User berhasil dibuat. Silakan verifikasi dengan kode OTP.", map[string]interface{}{
		"user":    result,
		"otp":     otp, // Remove this in production!
		"message": "Kode OTP telah dikirim ke nomor telepon Anda",
	})

	ctx.JSON(http.StatusCreated, response)
}

// VerifyOTP - Verify OTP and activate user account
func (c *authController) VerifyOTP(ctx *gin.Context) {
	var req request.VerifyOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Get stored OTP from Redis
	storedHash, err := c.otpHelper.GetOTP("register", req.PhoneNumber)
	if err != nil {
		res := helper.BuildErrorResponse("Kode OTP tidak valid atau sudah kadaluarsa", "INVALID_OTP", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Verify OTP
	if storedHash != helper.HashOTP(req.OTPCode) {
		res := helper.BuildErrorResponse("Kode OTP salah", "INVALID_OTP", "body", "OTP does not match", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Activate user account
	user, err := c.userService.VerifyAndActivate(req.PhoneNumber)
	if err != nil {
		res := helper.BuildErrorResponse("Gagal mengaktifkan akun", "ACTIVATION_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	// Delete OTP from Redis
	c.otpHelper.DeleteOTP("register", req.PhoneNumber)

	// Generate JWT token
	expiryTime := time.Now().Add(24 * time.Hour)
	token := c.jwtService.GenerateToken(user.Id, user.PhoneNumber, user.Email, expiryTime)

	response := helper.BuildResponse(true, "Akun berhasil diverifikasi dan diaktifkan", map[string]interface{}{
		"user":  user,
		"token": token,
	})

	ctx.JSON(http.StatusOK, response)
}

// Login - User login with phone/email and password
func (c *authController) Login(ctx *gin.Context) {
	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Login
	user, err := c.userService.Login(req)
	if err != nil {
		res := helper.BuildErrorResponse("Login gagal", "LOGIN_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	// Check if user is verified
	if !user.IsVerified {
		res := helper.BuildErrorResponse("Akun belum diverifikasi", "NOT_VERIFIED", "body", "Silakan verifikasi akun Anda terlebih dahulu", nil)
		ctx.JSON(http.StatusForbidden, res)
		return
	}

	// Check if user is active
	if !user.IsActive {
		res := helper.BuildErrorResponse("Akun tidak aktif", "ACCOUNT_INACTIVE", "body", "Akun Anda telah dinonaktifkan", nil)
		ctx.JSON(http.StatusForbidden, res)
		return
	}

	// Generate JWT token
	expiryTime := time.Now().Add(24 * time.Hour)
	token := c.jwtService.GenerateToken(user.Id, user.PhoneNumber, user.Email, expiryTime)

	response := helper.BuildResponse(true, "Login berhasil", map[string]interface{}{
		"user":  user,
		"token": token,
	})

	ctx.JSON(http.StatusOK, response)
}

// ResendOTP - Resend OTP for verification
func (c *authController) ResendOTP(ctx *gin.Context) {
	var req request.ResendOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := helper.BuildErrorResponse("Permintaan tidak valid", "INVALID_REQUEST", "body", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Check if user exists and not verified
	user, err := c.userService.FindByPhone(req.PhoneNumber)
	if err != nil {
		res := helper.BuildErrorResponse("User tidak ditemukan", "USER_NOT_FOUND", "body", err.Error(), nil)
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	if user.IsVerified {
		res := helper.BuildErrorResponse("Akun sudah diverifikasi", "ALREADY_VERIFIED", "body", "Akun Anda sudah diverifikasi", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	// Check rate limit
	if err := c.otpHelper.AllowRequest(req.PhoneNumber, 5, 15*time.Minute); err != nil {
		res := helper.BuildErrorResponse("Terlalu banyak permintaan", "RATE_LIMIT", "body", err.Error(), nil)
		ctx.JSON(http.StatusTooManyRequests, res)
		return
	}

	// Generate new OTP
	otp := helper.GenerateOTPCode(6)
	hashedOTP := helper.HashOTP(otp)

	// Save OTP to Redis
	if err := c.otpHelper.SaveOTP("register", req.PhoneNumber, hashedOTP, 15*time.Minute); err != nil {
		res := helper.BuildErrorResponse("Gagal mengirim OTP", "OTP_SEND_FAILED", "body", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	response := helper.BuildResponse(true, "Kode OTP baru telah dikirim", map[string]interface{}{
		"otp": otp, // Remove this in production!
	})

	ctx.JSON(http.StatusOK, response)
}
