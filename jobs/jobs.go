// Package jobs runs background work with River (riverqueue.io),
// a Postgres-backed job queue. Auth emails are enqueued here so HTTP
// handlers return fast and SMTP outages never fail API requests.
//
// River uses its own pgx pool alongside GORM and manages its schema
// (river_job, river_queue, ...) via rivermigrate in Setup.
package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"grip/models"
	"grip/repository"
	"grip/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
)

// QueueDefaultWorkers caps concurrent email sends.
const QueueDefaultWorkers = 10

// CompletedJobRetention is how long finished (completed/discarded) River
// jobs are kept for inspection before the daily purge deletes them.
const CompletedJobRetention = 7 * 24 * time.Hour

var (
	// Pool is River's dedicated pgx pool (separate from GORM's).
	Pool *pgxpool.Pool
	// Client is the shared River client used for inserts and working jobs.
	Client *river.Client[pgx.Tx]
)

// SendEmailArgs is the payload for a single outbound email.
// It carries everything the worker needs so delivery never
// depends on later database state.
type SendEmailArgs struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	HTML    string `json:"html,omitempty"`
}

// Kind implements river.JobArgs.
func (SendEmailArgs) Kind() string { return "send_email" }

// SendEmailWorker delivers emails enqueued as SendEmailArgs.
type SendEmailWorker struct {
	river.WorkerDefaults[SendEmailArgs]
}

// Work implements river.Worker.
func (w *SendEmailWorker) Work(ctx context.Context, job *river.Job[SendEmailArgs]) error {
	return utils.SendMail(job.Args.To, job.Args.Subject, job.Args.Body, job.Args.HTML)
}

// NotifyArgs is the payload for one inbox notification.
type NotifyArgs struct {
	UserID string         `json:"user_id"`
	Type   string         `json:"type"`
	Title  string         `json:"title"`
	Body   string         `json:"body"`
	Data   map[string]any `json:"data,omitempty"`
}

// Kind implements river.JobArgs.
func (NotifyArgs) Kind() string { return "notify" }

// NotifyWorker writes inbox rows for NotifyArgs jobs.
type NotifyWorker struct {
	river.WorkerDefaults[NotifyArgs]
}

// Work implements river.Worker.
func (w *NotifyWorker) Work(ctx context.Context, job *river.Job[NotifyArgs]) error {
	a := job.Args
	if a.UserID == "" || a.Title == "" {
		return fmt.Errorf("jobs: notify missing user_id/title")
	}
	return repository.CreateNotification(&models.Notification{
		ID:     uuid.New().String(),
		UserID: a.UserID,
		Type:   a.Type,
		Title:  a.Title,
		Body:   a.Body,
		Data:   a.Data,
	})
}

// EnqueueNotify queues an inbox notification (best-effort: callers log
// failures and continue, mirroring email enqueue semantics).
func EnqueueNotify(ctx context.Context, userID, typ, title, body string, data map[string]any) error {
	if Client == nil {
		return fmt.Errorf("jobs: client not started (call jobs.Setup first)")
	}
	_, err := Client.Insert(ctx, NotifyArgs{
		UserID: userID, Type: typ, Title: title, Body: body, Data: data,
	}, nil)
	if err != nil {
		return fmt.Errorf("jobs: insert notify: %w", err)
	}
	return nil
}

// CleanupArgs triggers the daily purge of old finished River jobs.
type CleanupArgs struct{}

// Kind implements river.JobArgs.
func (CleanupArgs) Kind() string { return "purge_completed_jobs" }

// CleanupWorker deletes completed/discarded jobs older than
// CompletedJobRetention so river_job can't grow unboundedly.
type CleanupWorker struct {
	river.WorkerDefaults[CleanupArgs]
}

// Work implements river.Worker.
func (w *CleanupWorker) Work(ctx context.Context, job *river.Job[CleanupArgs]) error {
	if Pool == nil {
		return fmt.Errorf("jobs: pool not started")
	}
	cutoff := time.Now().Add(-CompletedJobRetention)
	tag, err := Pool.Exec(ctx,
		`DELETE FROM river_job WHERE finalized_at < $1 AND state IN ('completed', 'discarded')`,
		cutoff)
	if err != nil {
		return fmt.Errorf("jobs: purge completed: %w", err)
	}
	slog.Info("Purged old river jobs", "rows", tag.RowsAffected())
	return nil
}

// Setup connects River, migrates its schema, registers workers and
// starts working jobs. Call once at boot (and in integration tests
// with the test database URL).
func Setup(ctx context.Context, databaseURL string) error {
	if databaseURL == "" {
		return fmt.Errorf("jobs: empty database URL")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("jobs: create pool: %w", err)
	}

	migrator, err := rivermigrate.New(riverpgxv5.New(pool), nil)
	if err != nil {
		pool.Close()
		return fmt.Errorf("jobs: migrator: %w", err)
	}
	if _, err := migrator.Migrate(ctx, rivermigrate.DirectionUp, nil); err != nil {
		pool.Close()
		return fmt.Errorf("jobs: migrate: %w", err)
	}

	workers := river.NewWorkers()
	river.AddWorker(workers, &SendEmailWorker{})
	river.AddWorker(workers, &NotifyWorker{})
	river.AddWorker(workers, &CleanupWorker{})

	client, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: QueueDefaultWorkers},
		},
		Workers: workers,
	})
	if err != nil {
		pool.Close()
		return fmt.Errorf("jobs: new client: %w", err)
	}
	if err := client.Start(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("jobs: start: %w", err)
	}

	// Daily purge of old finished jobs (unique ID: a second replica
	// registering the same schedule is a no-op, not a duplicate).
	client.PeriodicJobs().Add(river.NewPeriodicJob(
		river.PeriodicInterval(24*time.Hour),
		func() (river.JobArgs, *river.InsertOpts) { return CleanupArgs{}, nil },
		&river.PeriodicJobOpts{ID: "purge-completed-jobs"},
	))

	Pool = pool
	Client = client
	slog.Info("River client started")
	return nil
}

// Shutdown stops working jobs and closes the pool.
func Shutdown(ctx context.Context) error {
	if Client == nil {
		return nil
	}
	if err := Client.Stop(ctx); err != nil {
		return fmt.Errorf("jobs: stop: %w", err)
	}
	if Pool != nil {
		Pool.Close()
	}
	Client = nil
	Pool = nil
	return nil
}

// EnqueueSendEmail inserts a send_email job on the default queue.
func EnqueueSendEmail(ctx context.Context, to, subject, body, html string) error {
	if Client == nil {
		return fmt.Errorf("jobs: client not started (call jobs.Setup first)")
	}
	_, err := Client.Insert(ctx, SendEmailArgs{To: to, Subject: subject, Body: body, HTML: html}, nil)
	if err != nil {
		return fmt.Errorf("jobs: insert send_email: %w", err)
	}
	return nil
}

// EnqueueVerificationEmail queues the verify-your-email message for auth,
// containing both the OTP code and the verification link.
func EnqueueVerificationEmail(ctx context.Context, to, name, token, otp string) error {
	subject, body, html, err := utils.BuildVerificationEmail(name, token, otp)
	if err != nil {
		return fmt.Errorf("jobs: build verification email: %w", err)
	}
	return EnqueueSendEmail(ctx, to, subject, body, html)
}

// EnqueuePasswordResetEmail queues the password-reset message for auth.
func EnqueuePasswordResetEmail(ctx context.Context, to, name, token string) error {
	subject, body, html, err := utils.BuildPasswordResetEmail(name, token)
	if err != nil {
		return fmt.Errorf("jobs: build reset email: %w", err)
	}
	return EnqueueSendEmail(ctx, to, subject, body, html)
}
