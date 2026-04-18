package routes

import "github.com/gin-gonic/gin"

func SetupPublicRoutes(rg *gin.RouterGroup, opts any) {
	//1. Handler definition

	//2. Middleware for public routes can be added here (if needed)

	//3. Route definition
	// Example public route
	rg.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the public API endpoint",
		})
	})
}