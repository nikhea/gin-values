package utils

import (
	"fmt"
	"strings"
)

// Redis key namespace. Every key starts with the app prefix so one Redis
// instance can be shared safely with other services/databases.
const (
	KeyPrefix = "grip"

	KeyUsers    = "users"
	KeyProfiles = "profiles"
	KeyContacts = "contacts"
	KeyAuth     = "auth"
)

// Key joins namespace segments into a stable colon-separated key:
// Key("users", id) -> "grip:users:<id>".
func Key(segments ...string) string {
	return KeyPrefix + ":" + strings.Join(segments, ":")
}

// UserKey addresses a cached user entity.
func UserKey(id string) string { return Key(KeyUsers, id) }

// UserEmailKey maps an email to a user ID (login fast-path).
func UserEmailKey(email string) string {
	return Key(KeyUsers, "email", strings.ToLower(strings.TrimSpace(email)))
}

// ProfileKey addresses a cached profile by owner user ID.
func ProfileKey(userID string) string { return Key(KeyProfiles, userID) }

// ContactKey addresses a cached contact entity.
func ContactKey(id string) string { return Key(KeyContacts, id) }

// ContactListKey addresses a cached contact-list page. Filters are part
// of the key so different queries never collide.
func ContactListKey(userID, contactType, search string, page, pageSize int) string {
	return Key(KeyContacts, "list",
		fmt.Sprintf("u=%s", userID),
		fmt.Sprintf("t=%s", contactType),
		fmt.Sprintf("q=%s", search),
		fmt.Sprintf("p=%d", page),
		fmt.Sprintf("s=%d", pageSize),
	)
}

// OTPAttemptsKey counts OTP guesses per email (alternative to the DB counter).
func OTPAttemptsKey(email string) string {
	return Key(KeyAuth, "otp", "attempts", strings.ToLower(strings.TrimSpace(email)))
}

// MatchUsers / MatchContacts / MatchAuth are SCAN patterns for bulk
// invalidation (e.g. after a deploy that changes cached shapes).
func MatchUsers() string    { return Key(KeyUsers, "*") }
func MatchProfiles() string { return Key(KeyProfiles, "*") }
func MatchContacts() string { return Key(KeyContacts, "*") }
func MatchAuth() string     { return Key(KeyAuth, "*") }
