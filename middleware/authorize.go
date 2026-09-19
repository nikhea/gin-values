package middleware

import (
	"log/slog"
	"net/http"

	"grip/authz"

	"github.com/gin-gonic/gin"
)

// Context key carrying the resolved organization ID.
const ContextOrgID = "orgID"

// OrgIDFromRequest resolves the organization scope from the X-Org-ID
// header (preferred) or the ?org_id= query param.
func OrgIDFromRequest(c *gin.Context) string {
	if orgID := c.GetHeader("X-Org-ID"); orgID != "" {
		return orgID
	}
	return c.Query("org_id")
}

// GetOrgID returns the organization ID stored by RequireOrgAccess.
func GetOrgID(c *gin.Context) (string, bool) {
	v, exists := c.Get(ContextOrgID)
	if !exists {
		return "", false
	}
	id, ok := v.(string)
	if !ok || id == "" {
		return "", false
	}
	return id, true
}

// RequireOrgAccess enforces Casbin RBAC for one resource/action inside the
// request's organization. Must run after AuthRequired (needs GetUserID).
// Missing org scope is 400; denial is 403; enforcer errors fail closed.
func RequireOrgAccess(obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := GetUserID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}
		orgID := OrgIDFromRequest(c)
		if orgID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "X-Org-ID header (or org_id query param) is required",
			})
			return
		}

		allowed, err := authz.Enforce(userID, orgID, obj, act)
		if err != nil {
			slog.Error("authz enforce failed", "error", err, "user", userID, "org", orgID)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
			})
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
			})
			return
		}

		c.Set(ContextOrgID, orgID)
		c.Next()
	}
}
