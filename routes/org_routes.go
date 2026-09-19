package routes

import (
	"grip/authz"
	"grip/handlers"
	"grip/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterOrgRoutes mounts organization endpoints on the protected group.
// Creation/listing need only authentication; member management and org
// deletion are additionally gated by Casbin RBAC in the request's org.
func RegisterOrgRoutes(api *gin.RouterGroup) {
	orgs := api.Group("/orgs")
	{
		orgs.POST("/", handlers.CreateOrg)
		orgs.GET("/", handlers.ListOrgs)

		// The org is identified by :id here, so seed the scope the
		// authorizer reads instead of demanding a redundant header.
		scoped := orgs.Group("/:id", orgScopeFromPath)

		members := scoped.Group("/members",
			middleware.RequireOrgAccess(authz.ResourceMembers, authz.ActionWrite))
		{
			members.POST("", handlers.AddMember)
			members.PUT("/:uid", handlers.UpdateMember)
			members.DELETE("/:uid", handlers.RemoveMember)
		}

		scoped.DELETE("",
			middleware.RequireOrgAccess(authz.ResourceOrgs, authz.ActionDelete),
			handlers.DeleteOrg)
	}
}

// orgScopeFromPath copies the :id path param into X-Org-ID when the caller
// didn't scope the request explicitly.
func orgScopeFromPath(c *gin.Context) {
	if c.GetHeader("X-Org-ID") == "" && c.Query("org_id") == "" {
		c.Request.Header.Set("X-Org-ID", c.Param("id"))
	}
	c.Next()
}
