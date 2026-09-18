package tests

import (
	"net/http"
	"testing"
	"time"

	"grip/config"
	"grip/dto"
	"grip/repository"
	"grip/routes"
	services "grip/service"
	"grip/utils"
)

// setKnownOTP plants a deterministic code for tests (Register generates
// a random one that only travels by email).
func setKnownOTP(t *testing.T, email, code string, expiresAt time.Time) {
	t.Helper()
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		t.Fatalf("load user: %v", err)
	}
	user.OTPHash = utils.HashOTP(code)
	user.OTPExpiresAt = expiresAt
	user.OTPAttempts = 0
	if err := config.DB.Save(user).Error; err != nil {
		t.Fatalf("plant OTP: %v", err)
	}
}

func TestVerifyOTPFlow(t *testing.T) {
	requireTestDB(t)

	if _, err := services.Register(dto.RegisterRequest{
		Name: "OTP User", Email: "otp@example.com", Password: "supersecret123", Age: 30,
	}); err != nil {
		t.Fatalf("register: %v", err)
	}

	// Register stores only a hash with a future expiry.
	stored, err := repository.GetUserByEmail("otp@example.com")
	if err != nil {
		t.Fatalf("load user: %v", err)
	}
	if stored.OTPHash == "" || time.Until(stored.OTPExpiresAt) <= 0 {
		t.Fatal("expected OTP hash and future expiry after register")
	}

	setKnownOTP(t, "otp@example.com", "123456", time.Now().Add(10*time.Minute))

	// Wrong code fails and records an attempt.
	if _, err := services.VerifyOTP("otp@example.com", "000000"); err != services.ErrInvalidOTP {
		t.Fatalf("wrong code err = %v, want ErrInvalidOTP", err)
	}
	afterFail, _ := repository.GetUserByEmail("otp@example.com")
	if afterFail.OTPAttempts != 1 {
		t.Fatalf("attempts = %d, want 1", afterFail.OTPAttempts)
	}

	// Unknown email does not reveal anything.
	if _, err := services.VerifyOTP("nobody@example.com", "123456"); err != services.ErrInvalidOTP {
		t.Fatalf("unknown email err = %v, want ErrInvalidOTP", err)
	}

	// Correct code verifies and logs in.
	user, err := services.VerifyOTP("otp@example.com", "123456")
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !user.EmailVerified {
		t.Fatal("expected email verified")
	}
	if _, _, err := services.Login(dto.LoginRequest{Email: "otp@example.com", Password: "supersecret123"}); err != nil {
		t.Fatalf("login after OTP verify: %v", err)
	}

	// Re-verifying an already-verified email.
	if _, err := services.VerifyOTP("otp@example.com", "123456"); err != services.ErrAlreadyVerified {
		t.Fatalf("re-verify err = %v, want ErrAlreadyVerified", err)
	}
}

func TestVerifyOTPExpiryAndLockout(t *testing.T) {
	requireTestDB(t)

	if _, err := services.Register(dto.RegisterRequest{
		Name: "OTP Lock", Email: "otplock@example.com", Password: "supersecret123", Age: 30,
	}); err != nil {
		t.Fatalf("register: %v", err)
	}

	// Expired code.
	setKnownOTP(t, "otplock@example.com", "123456", time.Now().Add(-time.Minute))
	if _, err := services.VerifyOTP("otplock@example.com", "123456"); err != services.ErrInvalidOTP {
		t.Fatalf("expired err = %v, want ErrInvalidOTP", err)
	}

	// Fresh code, burn all attempts.
	setKnownOTP(t, "otplock@example.com", "123456", time.Now().Add(10*time.Minute))
	for i := 0; i < utils.OTPMaxAttempts-1; i++ {
		if _, err := services.VerifyOTP("otplock@example.com", "000000"); err != services.ErrInvalidOTP {
			t.Fatalf("attempt %d err = %v, want ErrInvalidOTP", i, err)
		}
	}
	if _, err := services.VerifyOTP("otplock@example.com", "000000"); err != services.ErrTooManyAttempts {
		t.Fatalf("final attempt err = %v, want ErrTooManyAttempts", err)
	}

	// Even the correct code is locked out now.
	if _, err := services.VerifyOTP("otplock@example.com", "123456"); err != services.ErrTooManyAttempts {
		t.Fatalf("locked correct err = %v, want ErrTooManyAttempts", err)
	}

	// Resending rotates the code and resets attempts.
	if err := services.ResendVerification("otplock@example.com"); err != nil {
		t.Fatalf("resend: %v", err)
	}
	setKnownOTP(t, "otplock@example.com", "654321", time.Now().Add(10*time.Minute))
	if _, err := services.VerifyOTP("otplock@example.com", "654321"); err != nil {
		t.Fatalf("verify after resend: %v", err)
	}
}

func TestVerifyOTPHTTP(t *testing.T) {
	requireTestDB(t)
	router := routes.Setup()

	if _, err := services.Register(dto.RegisterRequest{
		Name: "OTP HTTP", Email: "otphttp@example.com", Password: "supersecret123", Age: 30,
	}); err != nil {
		t.Fatalf("register: %v", err)
	}
	setKnownOTP(t, "otphttp@example.com", "123456", time.Now().Add(10*time.Minute))

	w := doRequest(t, router, "POST", "/api/auth/verify-otp",
		dto.VerifyOTPRequest{Email: "otphttp@example.com", Code: "000000"}, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("wrong code status = %d, want 400", w.Code)
	}

	w = doRequest(t, router, "POST", "/api/auth/verify-otp",
		dto.VerifyOTPRequest{Email: "otphttp@example.com", Code: "123456"}, "")
	if w.Code != http.StatusOK {
		t.Fatalf("correct code status = %d, body %s", w.Code, w.Body.String())
	}
}
