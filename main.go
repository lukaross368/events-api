package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukaross368/events-api/db"
	"github.com/lukaross368/events-api/processes"
	"github.com/lukaross368/events-api/routes"
)

const cleanupInterval = 1 * time.Minute

func main() {

	db.InitDB()

	go func() {
		for {
			time.Sleep(cleanupInterval)
			processes.CleanUpOldLimiters()
		}
	}()

	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run("0.0.0.0:8080")
}
