package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"grip/dto"
	"grip/jobs"
	"grip/models"
	"grip/repository"
	"grip/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Sentinel errors so handlers can map auth failures to status codes
// without leaking internal details to clients.
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already registered")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrAlreadyVerified    = errors.New("email already verified")
	ErrInvalidOTP         = errors.New("invalid or expired code")
	ErrTooManyAttempts    = errors.New("too many attempts, request a new code")
)

// Register creates a user with a bcrypt password hash, stores an email
// verification token, and enqueues the verification email as a River job
// (logged when SMTP is unconfigured so registration never fails for
// missing mail creds).
func Register(req dto.RegisterRequest) (*models.User, error) {
	if _, err := repository.GetUserByEmail(req.Email); err == nil {
		return nil, ErrEmailTaken
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	verifyToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		return nil, err
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:                uuid.New().String(),
		Name:              req.Name,
		Email:             req.Email,
		Age:               req.Age,
		PasswordHash:      hash,
		VerificationToken: verifyToken,
		OTPHash:           utils.HashOTP(otp),
		OTPExpiresAt:      time.Now().Add(utils.OTPExpiry),
	}

	if err := repository.CreateUser(user); err != nil {
		return nil, err
	}

	// Email delivery is async: enqueue failure must not fail registration.
	if err := jobs.EnqueueVerificationEmail(context.Background(), user.Email, user.Name, verifyToken, otp); err != nil {
		slog.Warn("verification email job failed", "error", err, "email", user.Email)
	}

	return user, nil
}

// VerifyEmail marks the user matching the verification token as verified.
func VerifyEmail(token string) (*models.User, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}
	user, err := repository.GetUserByVerificationToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}
	if user.EmailVerified {
		return nil, ErrAlreadyVerified
	}

	user.EmailVerified = true
	user.VerificationToken = ""
	if err := repository.UpdateUser(user); err != nil {
		return nil, err
	}
	_ = utils.Del(context.Background(), utils.UserKey(user.ID))
	return user, nil
}

// VerifyOTP verifies a user's email with the 6-digit code from the
// verification email. Codes expire after 10 minutes and allow at most
// 5 attempts before a fresh code must be requested.
func VerifyOTP(email, code string) (*models.User, error) {
	if email == "" || code == "" {
		return nil, ErrInvalidOTP
	}
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidOTP
		}
		return nil, err
	}
	if user.EmailVerified {
		return nil, ErrAlreadyVerified
	}
	if user.OTPHash == "" || time.Now().After(user.OTPExpiresAt) {
		return nil, ErrInvalidOTP
	}
	if user.OTPAttempts >= utils.OTPMaxAttempts {
		return nil, ErrTooManyAttempts
	}
	if !utils.CheckOTP(code, user.OTPHash) {
		user.OTPAttempts++
		if err := repository.UpdateUser(user); err != nil {
			return nil, err
		}
		if user.OTPAttempts >= utils.OTPMaxAttempts {
			return nil, ErrTooManyAttempts
		}
		return nil, ErrInvalidOTP
	}

	user.EmailVerified = true
	user.VerificationToken = ""
	user.OTPHash = ""
	user.OTPExpiresAt = time.Time{}
	user.OTPAttempts = 0
	if err := repository.UpdateUser(user); err != nil {
		return nil, err
	}
	_ = utils.Del(context.Background(), utils.UserKey(user.ID))
	return user, nil
}

// ResendVerification issues a fresh verification token and OTP code,
// then re-sends the verification email.
func ResendVerification(email string) error {
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal whether the email exists.
			return nil
		}
		return err
	}
	if user.EmailVerified {
		return ErrAlreadyVerified
	}

	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		return err
	}
	otp, err := utils.GenerateOTP()
	if err != nil {
		return err
	}
	user.VerificationToken = token
	user.OTPHash = utils.HashOTP(otp)
	user.OTPExpiresAt = time.Now().Add(utils.OTPExpiry)
	user.OTPAttempts = 0
	if err := repository.UpdateUser(user); err != nil {
		return err
	}
	// Never fail the request on mail errors; delivery is a background job.
	if err := jobs.EnqueueVerificationEmail(context.Background(), user.Email, user.Name, token, otp); err != nil {
		slog.Warn("resend verification email job failed", "error", err, "email", user.Email)
	}
	return nil
}

// Login validates credentials and email verification, then issues a JWT.
func Login(req dto.LoginRequest) (*models.User, string, string, error) {
	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", ErrInvalidCredentials
		}
		return nil, "", "", err
	}
	// Users created before auth (no password set) cannot log in.
	if user.PasswordHash == "" {
		return nil, "", "", ErrInvalidCredentials
	}
	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, "", "", ErrInvalidCredentials
	}
	if !user.EmailVerified {
		return nil, "", "", ErrEmailNotVerified
	}

	access, refresh, err := IssueTokenPair(user)
	if err != nil {
		return nil, "", "", err
	}
	return user, access, refresh, nil
}

// ForgotPassword creates a 1-hour reset token and emails it. Unknown
// emails return success to prevent account enumeration.
func ForgotPassword(email string) error {
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		return err
	}
	user.ResetToken = token
	user.ResetExpiresAt = time.Now().Add(time.Hour)
	if err := repository.UpdateUser(user); err != nil {
		return err
	}
	// Always succeed: log enqueue failures instead of revealing anything.
	if err := jobs.EnqueuePasswordResetEmail(context.Background(), user.Email, user.Name, token); err != nil {
		slog.Warn("password reset email job failed", "error", err, "email", user.Email)
	}
	return nil
}

// ResetPassword consumes a valid, unexpired reset token and sets a new password.
func ResetPassword(token, newPassword string) (*models.User, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}
	user, err := repository.GetUserByResetToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}
	if user.ResetToken == "" || time.Now().After(user.ResetExpiresAt) {
		return nil, ErrInvalidToken
	}

	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = hash
	user.ResetToken = ""
	user.ResetExpiresAt = time.Time{}
	if err := repository.UpdateUser(user); err != nil {
		return nil, err
	}
	_ = utils.Del(context.Background(), utils.UserKey(user.ID))
	return user, nil
}
