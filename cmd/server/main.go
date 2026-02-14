package main

import (
	"kavania2002/puzzle-together-server/api"
	lib "kavania2002/puzzle-together-server/lib/ws"

	"github.com/gin-gonic/gin"
)

// main creates a WebSocket manager, registers API routes with a default Gin router, and starts the HTTP server.
// It initializes a new ws manager via lib.NewManager(), registers routes using api.RegisterRoutes(manager, r),
// and runs the Gin server on the default listen address.
func main() {
	manager := lib.NewManager()

	r := gin.Default()
	api.RegisterRoutes(manager, r)
	r.Run()
}