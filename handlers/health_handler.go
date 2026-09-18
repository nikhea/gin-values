package handlers

import (
	"net/http"

	"gin-learn/config"
	"gin-learn/dto"
	"gin-learn/jobs"

	"github.com/gin-gonic/gin"
)

// Compile-time reference so Swagger can resolve dto.MessageEnvelope.
var _ = dto.MessageEnvelope{}

// Health godoc
// @Summary Health check
// @Tags health
// @Produce json
// @Success 200 {object} dto.MessageEnvelope
// @Router /health [get]
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// Readyz godoc
// @Summary Readiness probe (database + job queue)
// @Tags health
// @Produce json
// @Success 200 {object} dto.MessageEnvelope
// @Failure 503 {object} dto.ErrorEnvelope
// @Router /readyz [get]
func Readyz(c *gin.Context) {
	if config.DB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "database not ready"})
		return
	}
	sqlDB, err := config.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "database not ready"})
		return
	}
	if err := sqlDB.PingContext(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "database not ready"})
		return
	}
	if jobs.Client == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "job queue not ready"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
