// Command migrate runs versioned DB migrations manually:
//
//	go run ./cmd/migrate up
//	go run ./cmd/migrate down
//	go run ./cmd/migrate version
package main

import (
	"fmt"
	"log"
	"os"

	"gin-learn/config"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config.LoadEnv()

	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./cmd/migrate [up|down|version|force <version>]")
	}

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

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("Already up to date")
				return
			}
			log.Fatal("Migrate up failed:", err)
		}
		fmt.Println("Migrations up applied")
	case "down":
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("Already at base")
				return
			}
			log.Fatal("Migrate down failed:", err)
		}
		fmt.Println("Migrations down applied")
	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				fmt.Println("No migrations applied yet")
				return
			}
			log.Fatal("Version failed:", err)
		}
		fmt.Printf("version=%d dirty=%v\n", v, dirty)
	default:
		log.Fatal("unknown command, use up|down|version")
	}
}
