package tests

import (
	"context"
	"testing"

	"gin-learn/jobs"

	"github.com/riverqueue/river"
)

func TestSendEmailArgsKind(t *testing.T) {
	var args jobs.SendEmailArgs
	if got := args.Kind(); got != "send_email" {
		t.Fatalf("Kind() = %q, want %q", got, "send_email")
	}
}

func TestSendEmailWorkerLogOnly(t *testing.T) {
	// No SMTP creds: delivery logs and succeeds, no DB needed.
	t.Setenv("EMAIL_ADDRESS", "")
	t.Setenv("EMAIL_PASSWORD", "")

	w := &jobs.SendEmailWorker{}
	job := &river.Job[jobs.SendEmailArgs]{Args: jobs.SendEmailArgs{
		To: "user@example.com", Subject: "Hi", Body: "hello", HTML: "<p>hello</p>",
	}}
	if err := w.Work(context.Background(), job); err != nil {
		t.Fatalf("Work: %v", err)
	}
}

func TestEnqueueWithoutSetupFails(t *testing.T) {
	old := jobs.Client
	jobs.Client = nil
	defer func() { jobs.Client = old }()

	if err := jobs.EnqueueSendEmail(context.Background(), "a@b.c", "s", "b", "<p>b</p>"); err == nil {
		t.Fatal("expected error when River client is not started")
	}
}
