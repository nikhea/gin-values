package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies all pending up migrations from the migrations/ dir.
// DATABASE_URL must be set (e.g. postgresql://user:password@localhost:5432/dbname?sslmode=disable).
func RunMigrations() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatal("Migrate init failed:", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Println("Migrate source close error:", srcErr)
		}
		if dbErr != nil {
			log.Println("Migrate db close error:", dbErr)
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("Migrations already up to date")
			return
		}
		log.Fatal("Migrate up failed:", err)
	}

	fmt.Println("Migrations applied")
}

// MigrationVersion prints the current migration version for debugging.
func MigrationVersion() (uint, bool) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatal("Migrate init failed:", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	v, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			fmt.Println("No migrations applied yet")
			return 0, false
		}
		log.Fatal("Migrate version failed:", err)
	}

	fmt.Printf("Migration version: %d (dirty: %v)\n", v, dirty)
	return v, dirty
}
