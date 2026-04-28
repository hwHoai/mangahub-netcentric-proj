package routes

import (
	"crypto/rsa"

	"mangahub/cmd/api-server/middleware"

	"github.com/gin-gonic/gin"
)

type PrivateRouteOpts struct {
	PublicKey *rsa.PublicKey
}

func SetupPrivateRoutes(rg *gin.RouterGroup, opts *PrivateRouteOpts) {
	authMiddleware := middleware.NewAuthMiddleware(opts.PublicKey)
	rg.Use(authMiddleware.Handler())
	
	// Route definition
	// Example protected route
	rg.GET("/secret", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "This is a secret message",
		})
	})
}