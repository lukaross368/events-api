package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lukaross368/events-api/middlewares"
)

func RegisterRoutes(server *gin.Engine) {
	root := server.Group("/")
	root.Use(middlewares.RateLimit)

	root.GET("/events", getEvents)
	root.GET("/events/:id", getEvent)

	root.POST("/signup", signup)
	root.POST("/login", login)

	authenticated := root.Group("/")
	authenticated.Use(middlewares.Authenticate)

	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)

	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.DELETE("/events/:id/register", cancelEvent)
}
