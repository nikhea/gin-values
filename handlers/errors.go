package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// internalError logs the real failure with request context and responds
// with a generic 500. Handler internals (SQL, driver, filesystem details)
// must never reach clients.
func internalError(c *gin.Context, err error) {
	requestID, _ := c.Get("requestID")
	slog.Error("internal error",
		"requestID", requestID,
		"method", c.Request.Method,
		"path", c.FullPath(),
		"error", err,
	)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
