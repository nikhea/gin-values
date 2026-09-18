package utils

import (
	"encoding/hex"
	"testing"
)

func TestGenerateSecureToken(t *testing.T) {
	a, err := GenerateSecureToken(32)
	if err != nil {
		t.Fatalf("GenerateSecureToken: %v", err)
	}
	if len(a) != 64 {
		t.Fatalf("len = %d, want 64 (32 bytes hex)", len(a))
	}
	if _, err := hex.DecodeString(a); err != nil {
		t.Fatalf("not valid hex: %v", err)
	}

	b, err := GenerateSecureToken(32)
	if err != nil {
		t.Fatalf("GenerateSecureToken: %v", err)
	}
	if a == b {
		t.Fatal("expected unique tokens")
	}
}
