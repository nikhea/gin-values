package tests

import (
	"testing"
	"unicode"

	"gin-learn/utils"
)

func TestGenerateOTPFormat(t *testing.T) {
	code, err := utils.GenerateOTP()
	if err != nil {
		t.Fatalf("GenerateOTP: %v", err)
	}
	if len(code) != utils.OTPLength {
		t.Fatalf("len = %d, want %d", len(code), utils.OTPLength)
	}
	for _, r := range code {
		if !unicode.IsDigit(r) {
			t.Fatalf("non-digit in OTP %q", code)
		}
	}
}

func TestGenerateOTPUnique(t *testing.T) {
	a, err := utils.GenerateOTP()
	if err != nil {
		t.Fatalf("GenerateOTP: %v", err)
	}
	b, err := utils.GenerateOTP()
	if err != nil {
		t.Fatalf("GenerateOTP: %v", err)
	}
	if a == b {
		t.Fatal("expected unique codes")
	}
}

func TestHashAndCheckOTP(t *testing.T) {
	hash := utils.HashOTP("482914")
	if hash == "482914" {
		t.Fatal("hash must not equal plaintext")
	}
	if !utils.CheckOTP("482914", hash) {
		t.Fatal("expected correct code to verify")
	}
	if utils.CheckOTP("000000", hash) {
		t.Fatal("expected wrong code to fail")
	}
	if utils.CheckOTP("", hash) || utils.CheckOTP("482914", "") {
		t.Fatal("expected empty code/hash to fail")
	}
}
