package main

import (
	"kavania2002/puzzle-together-server/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	api.RegisterRoutes(r)
	r.Run()
}
