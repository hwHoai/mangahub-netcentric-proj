package routes

import (
	"mangahub/cmd/api-server/controllers"
	"mangahub/proto/session"
	"mangahub/proto/user"

	"github.com/gin-gonic/gin"
)	

type PublicRouteOpts struct {
	gRPCUserClient    user.GRPCUserServiceClient
	gRPCSessionClient session.GRPCSessionServiceClient
}

func SetupPublicRoutes(rg *gin.RouterGroup, opts *PublicRouteOpts) {
	//1. Handler definition (if needed)
	grpcUserClient := opts.gRPCUserClient
	grpcSessionClient := opts.gRPCSessionClient

	authController := controllers.NewAuthController(grpcUserClient, grpcSessionClient)

	//2. Middleware for public routes can be added here (if needed)

	//3. Route definition
	// Example public route
	rg.POST("/login", authController.LoginByUsername)
	rg.POST("/signup", authController.SignupByUsername)
}
