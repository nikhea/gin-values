package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plaintext password with bcrypt (DefaultCost).
// Password length is capped by the caller's validation (bcrypt max 72 bytes).
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword compares a bcrypt hash against a plaintext password.
func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
