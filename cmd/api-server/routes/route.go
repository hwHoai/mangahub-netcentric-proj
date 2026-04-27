package routes

import (
	"mangahub/internal/grpc"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	//1. Define gRPC clients for services
	grpcUserClient, _, err := grpc.NewUserGRPCClient()
	if err != nil {
		panic(err)
	}

	grpcSessionClient, _, err := grpc.NewSessionGRPCClient()
	if err != nil {
		panic(err)
	}

	//2. Route definition
	v1 := r.Group("api/v1")
	
	// Health Check
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	private_route_opts := &PrivateRouteOpts{
		// Add any dependencies needed for private routes here
	}
	SetupPrivateRoutes(v1, private_route_opts)

	public_route_opts := &PublicRouteOpts{
		gRPCUserClient:    grpcUserClient,
		gRPCSessionClient: grpcSessionClient,
	}
	SetupPublicRoutes(v1, public_route_opts)

}
