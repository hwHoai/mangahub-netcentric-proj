package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, opts any) {
	//1. Handler definition (if needed)
	
	//2. Route definition
	v1 := r.Group("api/v1")
	
	// Health Check
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	private_route_opts := struct {
		// Add any options needed for private routes here
	}{
		// Initialize options if needed
	}
	SetupPrivateRoutes(v1, private_route_opts)

	public_route_opts := struct {
		// Add any options needed for public routes here
	}{
		// Initialize options if needed
	}
	SetupPublicRoutes(v1, public_route_opts)

}