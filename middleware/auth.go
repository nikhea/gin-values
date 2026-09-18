package middleware

import (
	"net/http"
	"strings"

	"grip/utils"

	"github.com/gin-gonic/gin"
)

// Context keys set by AuthRequired for downstream handlers.
const (
	ContextUserID    = "userID"
	ContextUserEmail = "userEmail"
)

// AuthRequired validates a JWT from the Authorization header:
//
//	Authorization: Bearer <token>
//
// A bare token without the Bearer scheme is also accepted for
// convenience (e.g. pasted straight from a login response).
//
// On success it stores the user ID/email in the gin context and calls Next.
// On failure it aborts with 401.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := bearerToken(c.GetHeader("Authorization"))
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or malformed Authorization header",
			})
			return
		}

		claims, err := utils.ParseToken(token)
		if err != nil || claims.Subject == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		c.Set(ContextUserID, claims.Subject)
		c.Set(ContextUserEmail, claims.Email)
		c.Next()
	}
}

// bearerToken extracts the JWT from an Authorization header value,
// accepting both "Bearer <token>" and a bare "<token>".
func bearerToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}
	if parts := strings.SplitN(header, " ", 2); len(parts) == 2 {
		if !strings.EqualFold(parts[0], "bearer") {
			return ""
		}
		return strings.TrimSpace(parts[1])
	}
	return header
}

// GetUserID returns the authenticated user's ID from the gin context.
// It is the Go equivalent of Express's req.userId: call it in any
// handler behind AuthRequired().
//
//	userID, ok := middleware.GetUserID(c)
//	if !ok { ... } // only possible if AuthRequired() was bypassed
func GetUserID(c *gin.Context) (string, bool) {
	v, exists := c.Get(ContextUserID)
	if !exists {
		return "", false
	}
	id, ok := v.(string)
	if !ok || id == "" {
		return "", false
	}
	return id, true
}

// GetUserEmail returns the authenticated user's email from the gin context.
func GetUserEmail(c *gin.Context) (string, bool) {
	v, exists := c.Get(ContextUserEmail)
	if !exists {
		return "", false
	}
	email, ok := v.(string)
	if !ok || email == "" {
		return "", false
	}
	return email, true
}
