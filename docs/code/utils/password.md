# `utils/password.go`

bcrypt password hashing (`golang.org/x/crypto/bcrypt`, `DefaultCost`).

## Functions

- `HashPassword(password) (string, error)` — salted hash for storage. Callers enforce the 72-byte bcrypt limit via `max=72` binding tags.
- `CheckPassword(hash, password) error` — comparison; an empty hash (pre-auth users) always fails, so legacy rows can never log in.

Unit-tested in `tests/password_test.go` (round-trip, wrong password, salting, empty hash).
