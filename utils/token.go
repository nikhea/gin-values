package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecureToken returns a hex-encoded random token (n bytes of entropy).
func GenerateSecureToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
