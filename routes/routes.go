package routes

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"grip/handlers"
	"grip/middleware"

	_ "grip/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup registers all API routes and returns the configured router.
//
// /api/auth/* is public (register, login, verify, password reset).
// Everything else requires a JWT via the Authorization: Bearer header.
func Setup() *gin.Engine {
	router := gin.Default()

	// Trust only explicitly configured proxies (TRUSTED_PROXIES, comma
	// separated IPs/CIDRs). Empty means trust none, so ClientIP falls back
	// to the direct peer — correct when exposed without a proxy, and keeps
	// rate limiting and audit IPs honest behind one when configured.
	if err := router.SetTrustedProxies(trustedProxies()); err != nil {
		slog.Error("Invalid TRUSTED_PROXIES", "error", err)
		os.Exit(1)
	}

	router.Use(middleware.RequestID())
	router.Use(corsMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Liveness probe (public, no auth): used by Docker healthchecks.
	router.GET("/health", handlers.Health)
	// Readiness probe (public): database + job queue must be up.
	router.GET("/readyz", handlers.Readyz)

	// Publicly served uploaded files (avatars): /uploads/avatars/...
	router.Static("/uploads", "./uploads")

	RegisterAuthRoutes(router)

	protected := router.Group("/api")
	protected.Use(middleware.RateLimit("api", middleware.APIRateLimit))
	protected.Use(middleware.AuthRequired())
	{
		users := protected.Group("/users")
		{
			RegisterUserRoutes(users)
			RegisterProfileRoutes(users)
		}

		RegisterContactRoutes(protected)
		RegisterOrgRoutes(protected)
		RegisterAuditRoutes(protected)
		RegisterNotificationRoutes(protected)
	}

	return router
}

// trustedProxies parses TRUSTED_PROXIES (comma-separated IPs/CIDRs).
// Empty/missing means trust no proxies.
func trustedProxies() []string {
	raw := strings.TrimSpace(os.Getenv("TRUSTED_PROXIES"))
	if raw == "" {
		return nil
	}
	var out []string
	for _, p := range strings.Split(raw, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// corsMiddleware allows only explicitly configured origins
// (CORS_ALLOWED_ORIGINS, comma-separated). Unset means no CORS headers
// are emitted (same-origin / non-browser clients unaffected).
func corsMiddleware() gin.HandlerFunc {
	raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if raw == "" {
		return func(c *gin.Context) { c.Next() }
	}
	var origins []string
	for _, o := range strings.Split(raw, ",") {
		if o = strings.TrimSpace(o); o != "" {
			origins = append(origins, o)
		}
	}
	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
