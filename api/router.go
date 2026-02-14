package api

import (
	"net/http"

	ws "kavania2002/puzzle-together-server/lib/ws"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers HTTP routes on the provided Gin engine.
// It mounts:
// - GET "/" which responds with "Hello World",
// - GET "/ws" which delegates WebSocket handling to the provided Manager,
// - an API group "/api/v1" containing GET "/health" served by healthHandler.
func RegisterRoutes(m *ws.Manager, r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello World")
	})

	// Websocket
	r.GET("/ws", m.ServeWS)

	apiGroup := r.Group("/api/v1")

	// HEALTH ROUTE
	apiGroup.GET("/health", healthHandler)
}
