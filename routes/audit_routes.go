package routes

import (
	"grip/authz"
	"grip/handlers"
	"grip/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAuditRoutes mounts the append-only audit trail on the protected
// group. Listing and reads require org admin+ (members:write) in the
// request's organization.
func RegisterAuditRoutes(api *gin.RouterGroup) {
	logs := api.Group("/audit-logs",
		middleware.RequireOrgAccess(authz.ResourceMembers, authz.ActionWrite))
	{
		logs.GET("/", handlers.ListAuditLogs)
		logs.GET("/:id", handlers.GetAuditLog)
	}
}
