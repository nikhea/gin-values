// Command migrate runs versioned DB migrations manually:
//
//	go run ./cmd/migrate up
//	go run ./cmd/migrate down
//	go run ./cmd/migrate version
package main

import (
	"fmt"
	"log/slog"
	"os"

	"errors"
	"grip/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func fatal(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}

func main() {
	config.InitLogger()
	config.LoadEnv()

	if len(os.Args) < 2 {
		slog.Error("usage: go run ./cmd/migrate [up|down|version|force <version>]")
		os.Exit(1)
	}

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

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("Already up to date")
				return
			}
			fatal("Migrate up failed", err)
		}
		fmt.Println("Migrations up applied")
	case "down":
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("Already at base")
				return
			}
			fatal("Migrate down failed", err)
		}
		fmt.Println("Migrations down applied")
	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				fmt.Println("No migrations applied yet")
				return
			}
			fatal("Version failed", err)
		}
		fmt.Printf("version=%d dirty=%v\n", v, dirty)
	default:
		slog.Error("unknown command, use up|down|version")
		os.Exit(1)
	}
}
