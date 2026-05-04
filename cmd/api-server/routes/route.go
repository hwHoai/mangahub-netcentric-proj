package routes

import (
	"crypto/rsa"
	"mangahub/pkg/clients"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) {
	//1. Define gRPC clients for services
	grpcUserClient, _, err := clients.NewUserGRPCClient()
	if err != nil {
		panic(err)
	}

	grpcSessionClient, _, err := clients.NewSessionGRPCClient()
	if err != nil {
		panic(err)
	}

	grpcMangaClient, _, err := clients.NewMangaGRPCClient()
	if err != nil {
		panic(err)
	}

	grpcUserMangaClient, _, err := clients.NewUserMangaGRPCClient()
	if err != nil {
		panic(err)
	}

	grpcChapterClient, _, err := clients.NewChapterGRPCClient()
	if err != nil {
		panic(err)
	}

	tcpChapterSyncClient := clients.NewTCPChapterSyncClient()

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
		GRPCMangaClient:   grpcMangaClient,
		GRPCChapterClient: grpcChapterClient,
		PrivateKey:        privateKey,
		PublicKey:         publicKey,
	}
	SetupPublicRoutes(v1, public_route_opts)

	private_route_opts := &PrivateRouteOpts{
		PublicKey:            publicKey,
		GRPCUserMangaClient:  grpcUserMangaClient,
		GRPCUserClient:       grpcUserClient,
		GRPCSessionClient:    grpcSessionClient,
		GRPCChapterClient:    grpcChapterClient,
		TCPChapterSyncClient: tcpChapterSyncClient,
	}
	SetupPrivateRoutes(v1, private_route_opts)
}