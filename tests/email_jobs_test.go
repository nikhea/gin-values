package tests

import (
	"strings"
	"testing"

	"grip/config"
	"grip/dto"
	"grip/jobs"
	services "grip/service"
)

func riverJobCount(t *testing.T, kind string) int64 {
	t.Helper()
	var n int64
	if err := config.DB.Table("river_job").Where("kind = ?", kind).Count(&n).Error; err != nil {
		t.Fatalf("count river jobs: %v", err)
	}
	return n
}

func latestRiverJobArgs(t *testing.T, kind string) string {
	t.Helper()
	var encoded string
	if err := config.DB.Table("river_job").Select("args").
		Where("kind = ?", kind).Order("id DESC").Limit(1).Scan(&encoded).Error; err != nil {
		t.Fatalf("load river job args: %v", err)
	}
	return encoded
}

func TestRegisterEnqueuesVerificationEmail(t *testing.T) {
	requireTestDB(t)
	if jobs.Client == nil {
		t.Fatal("expected River client to be started")
	}

	before := riverJobCount(t, "send_email")
	if _, err := services.Register(dto.RegisterRequest{
		Name: "River Test", Email: "river@example.com", Password: "supersecret123", Age: 30,
	}); err != nil {
		t.Fatalf("register: %v", err)
	}

	if got := riverJobCount(t, "send_email"); got != before+1 {
		t.Fatalf("send_email jobs = %d, want %d", got, before+1)
	}
	if args := latestRiverJobArgs(t, "send_email"); !strings.Contains(args, "river@example.com") {
		t.Fatalf("job args missing recipient: %s", args)
	}
}

func TestForgotPasswordEnqueuesResetEmail(t *testing.T) {
	requireTestDB(t)
	registerVerifiedUser(t, "River Reset", "river-reset@example.com", "oldpassword1")

	before := riverJobCount(t, "send_email")
	if err := services.ForgotPassword("river-reset@example.com"); err != nil {
		t.Fatalf("forgot: %v", err)
	}

	if got := riverJobCount(t, "send_email"); got != before+1 {
		t.Fatalf("send_email jobs = %d, want %d", got, before+1)
	}
	if args := latestRiverJobArgs(t, "send_email"); !strings.Contains(args, "river-reset@example.com") {
		t.Fatalf("job args missing recipient: %s", args)
	}
}
