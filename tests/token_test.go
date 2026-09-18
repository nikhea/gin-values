package tests

import (
	"encoding/hex"
	"testing"

	"gin-learn/utils"
)

func TestGenerateSecureToken(t *testing.T) {
	a, err := utils.GenerateSecureToken(32)
	if err != nil {
		t.Fatalf("GenerateSecureToken: %v", err)
	}
	if len(a) != 64 {
		t.Fatalf("len = %d, want 64 (32 bytes hex)", len(a))
	}
	if _, err := hex.DecodeString(a); err != nil {
		t.Fatalf("not valid hex: %v", err)
	}

	b, err := utils.GenerateSecureToken(32)
	if err != nil {
		t.Fatalf("GenerateSecureToken: %v", err)
	}
	if a == b {
		t.Fatal("expected unique tokens")
	}
}
