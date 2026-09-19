package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// SHA256Hex returns the hex digest of s. Used for storing opaque tokens
// (OTP codes, refresh tokens) so only hashes ever touch the database.
func SHA256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// GenerateSecureToken returns a hex-encoded random token (n bytes of entropy).
func GenerateSecureToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
