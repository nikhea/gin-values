package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"math/big"
	"time"
)

// One-time passcode settings for email verification.
const (
	// OTPLength is the number of digits in a verification code.
	OTPLength = 6
	// OTPExpiry is how long a code stays valid.
	OTPExpiry = 10 * time.Minute
	// OTPMaxAttempts caps guesses per code before a fresh one is required.
	OTPMaxAttempts = 5
)

// GenerateOTP returns a cryptographically random numeric code.
func GenerateOTP() (string, error) {
	digits := make([]byte, OTPLength)
	for i := range digits {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		digits[i] = byte('0' + n.Int64())
	}
	return string(digits), nil
}

// HashOTP returns the SHA-256 hex digest of a code. Only the digest
// is stored in the database; the plain code travels by email only.
func HashOTP(code string) string {
	return SHA256Hex(code)
}

// CheckOTP compares a candidate code against a stored digest in
// constant time.
func CheckOTP(code, hash string) bool {
	if code == "" || hash == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(HashOTP(code)), []byte(hash)) == 1
}
