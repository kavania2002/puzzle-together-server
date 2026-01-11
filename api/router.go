package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello World")
	})

	apiGroup := r.Group("/api/v1")

	// HEALTH ROUTE
	apiGroup.GET("/health", healthHandler)
}
