package routes

import (
	"crypto/rsa"
	"mangahub/internal/grpc"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) {
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
	public_route_opts := &PublicRouteOpts{
		gRPCUserClient:    grpcUserClient,
		gRPCSessionClient: grpcSessionClient,
		PrivateKey:        privateKey,
		PublicKey:         publicKey,
	}
	SetupPublicRoutes(v1, public_route_opts)

	private_route_opts := &PrivateRouteOpts{
		PublicKey: publicKey,
	}
	SetupPrivateRoutes(v1, private_route_opts)

	

}
