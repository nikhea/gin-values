package handlers

import (
	"net/http"

	"gin-learn/dto"

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
