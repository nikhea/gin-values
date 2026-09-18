package services

import (
	"errors"
	"log"
	"time"

	"gin-learn/dto"
	"gin-learn/models"
	"gin-learn/repository"
	"gin-learn/utils"

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
)

// Register creates a user with a bcrypt password hash, stores an email
// verification token, and sends the verification email (logged when SMTP
// is unconfigured so registration never fails for missing mail creds).
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

	user := &models.User{
		ID:                uuid.New().String(),
		Name:              req.Name,
		Email:             req.Email,
		Age:               req.Age,
		PasswordHash:      hash,
		VerificationToken: verifyToken,
	}

	if err := repository.CreateUser(user); err != nil {
		return nil, err
	}

	// Email failure must not fail registration; log and continue.
	if err := utils.SendVerificationEmail(user.Email, user.Name, verifyToken); err != nil {
		log.Println("verification email failed:", err)
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
	return user, nil
}

// ResendVerification issues a fresh verification token and re-sends it.
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
	user.VerificationToken = token
	if err := repository.UpdateUser(user); err != nil {
		return err
	}
	// Never fail the request on mail errors; the mailer logs when disabled.
	if err := utils.SendVerificationEmail(user.Email, user.Name, token); err != nil {
		log.Println("resend verification email failed:", err)
	}
	return nil
}

// Login validates credentials and email verification, then issues a JWT.
func Login(req dto.LoginRequest) (*models.User, string, error) {
	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}
	// Users created before auth (no password set) cannot log in.
	if user.PasswordHash == "" {
		return nil, "", ErrInvalidCredentials
	}
	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, "", ErrInvalidCredentials
	}
	if !user.EmailVerified {
		return nil, "", ErrEmailNotVerified
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
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
	// Always succeed: log mail failures instead of revealing anything.
	if err := utils.SendPasswordResetEmail(user.Email, user.Name, token); err != nil {
		log.Println("password reset email failed:", err)
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
	return user, nil
}
