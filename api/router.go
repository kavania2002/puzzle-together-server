package api

import (
	"net/http"

	lib "kavania2002/puzzle-together-server/lib/ws"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(m *lib.Manager, r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello World")
	})

	// Websocket
	r.GET("/ws", m.ServeWS)

	apiGroup := r.Group("/api/v1")

	// HEALTH ROUTE
	apiGroup.GET("/health", healthHandler)
}
