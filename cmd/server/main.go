package main

import (
	"kavania2002/puzzle-together-server/api"
	lib "kavania2002/puzzle-together-server/lib/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	manager := lib.NewManager()

	r := gin.Default()
	api.RegisterRoutes(manager, r)
	r.Run()
}
