package main

import (
	"kavania2002/puzzle-together-server/api"
	ws "kavania2002/puzzle-together-server/lib/ws"
	"log"

	"github.com/gin-gonic/gin"
)

// main creates a WebSocket manager, registers API routes with a default Gin router, and starts the HTTP server.
// It initializes a new ws manager via ws.NewManager(), registers routes using api.RegisterRoutes(manager, r),
// and runs the Gin server on the default listen address.
func main() {
	manager := ws.NewManager()

	r := gin.Default()
	api.RegisterRoutes(manager, r)
	if err := r.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
