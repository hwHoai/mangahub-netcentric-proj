package routes

import "github.com/gin-gonic/gin"

type PrivateRouteOpts struct {
	// Add any dependencies needed for private routes here (e.g., services, repositories)
}

func SetupPrivateRoutes(rg *gin.RouterGroup, opts *PrivateRouteOpts) {
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