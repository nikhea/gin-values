package utils

import (
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("supersecret123")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == "supersecret123" {
		t.Fatal("hash must not equal plaintext")
	}
	if err := CheckPassword(hash, "supersecret123"); err != nil {
		t.Fatalf("CheckPassword correct: %v", err)
	}
	if err := CheckPassword(hash, "wrongpassword"); err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestHashPasswordSalts(t *testing.T) {
	a, err := HashPassword("same-password")
	if err != nil {
		t.Fatalf("hash a: %v", err)
	}
	b, err := HashPassword("same-password")
	if err != nil {
		t.Fatalf("hash b: %v", err)
	}
	if a == b {
		t.Fatal("expected different hashes for same password (bcrypt salt)")
	}
}

func TestCheckPasswordEmptyHash(t *testing.T) {
	// Users created before auth have no password set; they must never authenticate.
	if err := CheckPassword("", "anything"); err == nil {
		t.Fatal("expected error for empty hash")
	}
}
