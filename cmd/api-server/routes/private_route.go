package routes

import "github.com/gin-gonic/gin"

func SetupPrivateRoutes(rg *gin.RouterGroup, opts any) {
	//1. Handler definition

	//2. Middleware for private routes can be added here (e.g., JWT Authentication)
	
	//3. Route definition
	// Example protected route
	rg.GET("/secret", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "This is a secret message",
		})
	})
}