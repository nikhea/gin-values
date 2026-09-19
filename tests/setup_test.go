// Package tests holds DB-backed integration tests for the API.
//
// They need Postgres. By default they use a dedicated test database
// (gin_app_test on localhost, same credentials as development).
// Override with TEST_DATABASE_URL. If no database is reachable the
// tests skip instead of failing, so `go test ./...` stays green in
// environments without Postgres.
package tests

import (
	"context"
	"os"
	"strings"
	"testing"

	"grip/authz"
	"grip/config"
	"grip/jobs"
	"grip/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testDatabaseURL(t *testing.T) string {
	t.Helper()
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		return url
	}
	return "postgresql://admin:adminpassword@localhost:5432/gin_app_test?sslmode=disable"
}

// requireTestDB connects to the test database (creating it if needed),
// points config.DB at it, migrates the schema and truncates all tables.
// It skips the test when Postgres is unreachable.
func requireTestDB(t *testing.T) {
	t.Helper()

	t.Setenv("JWT_SECRET", "test-secret-for-integration-tests")
	t.Setenv("JWT_TTL_HOURS", "24")
	// Disable rate limiting for the suite (all requests share one test IP).
	t.Setenv("RATE_LIMIT_AUTH_PER_MINUTE", "100000")
	t.Setenv("RATE_LIMIT_API_PER_MINUTE", "100000")
	// Never send real email from tests: the mailer logs and succeeds.
	t.Setenv("EMAIL_ADDRESS", "")
	t.Setenv("EMAIL_PASSWORD", "")
	t.Setenv("APP_URL", "http://localhost:8080")

	dsn := testDatabaseURL(t)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		if tryCreateTestDatabase(dsn) {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		}
	}
	if err != nil {
		t.Skipf("postgres not available (%v), skipping integration test", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Profile{}, &models.Contact{}, &models.RefreshToken{}, &models.Organization{}, &models.Membership{}, &models.AuditLog{}, &models.Notification{}); err != nil {
		t.Fatalf("migrate test schema: %v", err)
	}

	// Mirror migration 000008: AutoMigrate creates plain unique indexes,
	// but production uses partial ones so soft-deleted rows never block
	// reuse (e.g. re-registering an email). Recreate them here for parity.
	for _, stmt := range []string{
		`DROP INDEX IF EXISTS idx_users_email`,
		`CREATE UNIQUE INDEX idx_users_email ON users (email) WHERE deleted_at IS NULL`,
		`DROP INDEX IF EXISTS idx_profiles_user_id`,
		`CREATE UNIQUE INDEX idx_profiles_user_id ON profiles (user_id) WHERE deleted_at IS NULL`,
	} {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("partial unique index: %v", err)
		}
	}

	config.DB = db

	// Authorization backend shares the test database (casbin_rule table
	// is truncated with everything else in truncateAll).
	if err := authz.Init(db); err != nil {
		t.Fatalf("casbin setup: %v", err)
	}

	// Start River workers against the test database so service calls
	// that enqueue email jobs work (and are actually worked) in tests.
	// A previous test's client is stopped first; tests are sequential.
	_ = jobs.Shutdown(context.Background())
	if err := jobs.Setup(context.Background(), dsn); err != nil {
		t.Fatalf("river setup: %v", err)
	}

	truncateAll(t)

	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { truncateAll(t) })
	t.Cleanup(func() { _ = jobs.Shutdown(context.Background()) })
	// Uploaded test files land in tests/uploads (package working dir).
	t.Cleanup(func() { _ = os.RemoveAll("uploads") })
}

func truncateAll(t *testing.T) {
	t.Helper()
	if config.DB == nil {
		return
	}
	if err := config.DB.Exec("TRUNCATE users, profiles, contacts, organizations, memberships, casbin_rule, audit_logs, notifications, river_job RESTART IDENTITY CASCADE").Error; err != nil {
		t.Fatalf("truncate test tables: %v", err)
	}
}

// tryCreateTestDatabase creates the database named in dsn by connecting
// to the maintenance DB. Returns true if the subsequent connect is worth retrying.
func tryCreateTestDatabase(dsn string) bool {
	name := databaseName(dsn)
	if name == "" {
		return false
	}
	adminDSN := strings.Replace(dsn, "/"+name+"?", "/postgres?", 1)
	admin, err := gorm.Open(postgres.Open(adminDSN), &gorm.Config{})
	if err != nil {
		return false
	}
	sqlDB, err := admin.DB()
	if err != nil {
		return false
	}
	defer func() { _ = sqlDB.Close() }()
	// Ignore error: it already exists in the common retry case.
	_ = admin.Exec(`CREATE DATABASE "` + name + `"`).Error
	return true
}

func databaseName(dsn string) string {
	// postgresql://user:pass@host:port/name?params
	rest := dsn
	if i := strings.Index(rest, "://"); i >= 0 {
		rest = rest[i+3:]
	}
	i := strings.Index(rest, "/")
	if i < 0 {
		return ""
	}
	rest = rest[i+1:]
	if j := strings.Index(rest, "?"); j >= 0 {
		rest = rest[:j]
	}
	return rest
}
