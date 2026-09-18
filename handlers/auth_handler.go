package handlers

import (
	"errors"
	"net/http"

	"gin-learn/dto"
	"gin-learn/middleware"
	services "gin-learn/service"

	"github.com/gin-gonic/gin"
)

func writeAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrEmailTaken):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrEmailNotVerified):
		c.JSON(http.StatusForbidden, gin.H{
			"error":   err.Error(),
			"message": "Please verify your email, or resend the verification email",
		})
	case errors.Is(err, services.ErrInvalidToken):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrInvalidOTP):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrTooManyAttempts):
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrAlreadyVerified):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		internalError(c, err)
	}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration payload"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 409 {object} dto.ErrorEnvelope
// @Failure 500 {object} dto.ErrorEnvelope
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.Register(req)
	if err != nil {
		writeAuthError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registered. Please check your email to verify your account",
		"user":    user,
	})
}

// Login godoc
// @Summary Log in with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login payload"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 401 {object} dto.ErrorEnvelope
// @Failure 403 {object} dto.ErrorEnvelope
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := services.Login(req)
	if err != nil {
		writeAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in",
		"token":   token,
		"user":    user,
	})
}

// VerifyEmail godoc
// @Summary Verify email address
// @Tags auth
// @Produce json
// @Param token query string true "Verification token from email" example(1ef700f6a95416bb2a1d90a6ee3e8d62e78f9b3e084e405ab29fd722a57c96e8)
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorEnvelope
// @Router /auth/verify [get]
func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	user, err := services.VerifyEmail(token)
	if err != nil {
		writeAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified. You can now log in",
		"user":    user,
	})
}

// VerifyOTP godoc
// @Summary Verify email with OTP code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyOTPRequest true "Email and 6-digit code"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorEnvelope
// @Failure 429 {object} dto.ErrorEnvelope
// @Router /auth/verify-otp [post]
func VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.VerifyOTP(req.Email, req.Code)
	if err != nil {
		writeAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified. You can now log in",
		"user":    user,
	})
}

// ResendVerification godoc
// @Summary Resend verification email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ResendVerificationRequest true "Email payload"
// @Success 200 {object} dto.MessageEnvelope
// @Failure 400 {object} dto.ErrorEnvelope
// @Router /auth/resend-verification [post]
func ResendVerification(c *gin.Context) {
	var req dto.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ResendVerification(req.Email); err != nil {
		writeAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "If the email exists and is unverified, a new verification email was sent",
	})
}

// ForgotPassword godoc
// @Summary Request a password-reset email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Email payload"
// @Success 200 {object} dto.MessageEnvelope
// @Failure 400 {object} dto.ErrorEnvelope
// @Router /auth/forgot-password [post]
func ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.ForgotPassword(req.Email); err != nil {
		writeAuthError(c, err)
		return
	}

	// Always succeed to prevent account enumeration.
	c.JSON(http.StatusOK, gin.H{
		"message": "If the email exists, a password-reset email was sent",
	})
}

// ResetPassword godoc
// @Summary Reset password with a reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset payload"
// @Success 200 {object} dto.MessageEnvelope
// @Failure 400 {object} dto.ErrorEnvelope
// @Router /auth/reset-password [post]
func ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Allow the token via query (?token=...) for email-link convenience.
	if req.Token == "" {
		req.Token = c.Query("token")
	}
	if req.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	if _, err := services.ResetPassword(req.Token, req.NewPassword); err != nil {
		writeAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset. You can now log in",
	})
}

// Me godoc
// @Summary Get the authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.AuthResponse
// @Failure 401 {object} dto.ErrorEnvelope
// @Failure 404 {object} dto.MessageEnvelope
// @Router /auth/me [get]
func Me(c *gin.Context) {
	// req.userId equivalent: the ID AuthRequired() stored from the JWT.
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Authenticated",
		"user":    user,
	})
}
