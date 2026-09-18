package tests

import (
	"strings"
	"testing"

	"grip/utils"
)

func TestKeyBuilders(t *testing.T) {
	if got := utils.UserKey("abc"); got != "grip:users:abc" {
		t.Fatalf("UserKey = %q", got)
	}
	if got := utils.ProfileKey("u1"); got != "grip:profiles:u1" {
		t.Fatalf("ProfileKey = %q", got)
	}
	if got := utils.ContactKey("c1"); got != "grip:contacts:c1" {
		t.Fatalf("ContactKey = %q", got)
	}
	if got := utils.OTPAttemptsKey("User@Example.COM "); got != "grip:auth:otp:attempts:user@example.com" {
		t.Fatalf("OTPAttemptsKey = %q (want normalized)", got)
	}
	if got := utils.UserEmailKey("  A@B.c "); got != "grip:users:email:a@b.c" {
		t.Fatalf("UserEmailKey = %q (want normalized)", got)
	}

	a := utils.ContactListKey("u1", "email", "kai", 1, 10)
	b := utils.ContactListKey("u1", "email", "kai", 2, 10)
	c := utils.ContactListKey("u1", "email", "kai", 1, 10)
	if a == b {
		t.Fatal("different pages must produce different keys")
	}
	if a != c {
		t.Fatal("identical filters must produce identical keys")
	}
	if !strings.HasPrefix(a, "grip:contacts:list:") {
		t.Fatalf("ContactListKey = %q", a)
	}
}

func TestMatchPatterns(t *testing.T) {
	for _, tc := range []struct{ got, prefix string }{
		{utils.MatchUsers(), "grip:users:"},
		{utils.MatchProfiles(), "grip:profiles:"},
		{utils.MatchContacts(), "grip:contacts:"},
		{utils.MatchAuth(), "grip:auth:"},
	} {
		if !strings.HasPrefix(tc.got, tc.prefix) || !strings.HasSuffix(tc.got, "*") {
			t.Fatalf("pattern = %q, want %q*", tc.got, tc.prefix)
		}
	}
}
