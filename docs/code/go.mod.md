# `go.mod` / `go.sum`

Go module definition for `module grip` (Go 1.27.1). `go.sum` is the generated checksum lockfile — never edited by hand.

## Direct dependencies

| Module | Used for |
|---|---|
| `github.com/gin-gonic/gin` | HTTP router and middleware |
| `github.com/golang-jwt/jwt/v5` | Signing/parsing access tokens (`utils/jwt.go`) |
| `github.com/google/uuid` | UUID string IDs for users, profiles, contacts |
| `github.com/joho/godotenv` | Loading `.env` in `config/env.go` |
| `github.com/jackc/pgx/v5` | Postgres driver pool used by River (`jobs`) |
| `github.com/riverqueue/river` + `riverdriver/riverpgxv5` | Postgres-backed background job queue for emails |
| `github.com/golang-migrate/migrate/v4` | Versioned SQL migrations in `migrations/` |
| `github.com/swaggo/files`, `gin-swagger`, `swag` | Swagger UI endpoint and doc generation |
| `golang.org/x/crypto` | bcrypt password hashing (`utils/password.go`) |
| `gorm.io/driver/postgres`, `gorm.io/gorm` | ORM and Postgres dialect (`config.DB`) |

## Notes

- Unused requirements are pruned by `go mod tidy`, so a dependency only stays direct while some `.go` file imports it.
- Adding a dependency: `go get <module>@<version>` then `go mod tidy`.
