package config

import (
	"errors"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// fatal logs the message with the error and exits non-zero,
// mirroring log.Fatal now that slog is the default logger.
func fatal(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}

// RunMigrations applies all pending up migrations from the migrations/ dir.
// DATABASE_URL must be set (e.g. postgresql://user:password@localhost:5432/dbname?sslmode=disable).
func RunMigrations() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		slog.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		fatal("Migrate init failed", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			slog.Warn("Migrate source close error", "error", srcErr)
		}
		if dbErr != nil {
			slog.Warn("Migrate db close error", "error", dbErr)
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("Migrations already up to date")
			return
		}
		fatal("Migrate up failed", err)
	}

	slog.Info("Migrations applied")
}

// MigrationVersion prints the current migration version for debugging.
func MigrationVersion() (uint, bool) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		slog.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		fatal("Migrate init failed", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	v, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			slog.Info("No migrations applied yet")
			return 0, false
		}
		fatal("Migrate version failed", err)
	}

	slog.Info("Migration version", "version", v, "dirty", dirty)
	return v, dirty
}
