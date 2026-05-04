package routes

import (
	"crypto/rsa"
	"mangahub/cmd/api-server/controllers"
	"mangahub/proto/chapter"
	"mangahub/proto/manga"
	"mangahub/proto/session"
	"mangahub/proto/user"

	"github.com/gin-gonic/gin"
)

type PublicRouteOpts struct {
	gRPCUserClient    user.GRPCUserServiceClient
	gRPCSessionClient session.GRPCSessionServiceClient
	GRPCMangaClient   manga.GRPCMangaServiceClient
	GRPCChapterClient chapter.GRPCChapterServiceClient
	PrivateKey        *rsa.PrivateKey
	PublicKey         *rsa.PublicKey
}

func SetupPublicRoutes(rg *gin.RouterGroup, opts *PublicRouteOpts) {
	//1. Handler definition (if needed)
	grpcUserClient := opts.gRPCUserClient
	grpcSessionClient := opts.gRPCSessionClient
	grpcMangaClient := opts.GRPCMangaClient
	grpcChapterClient := opts.GRPCChapterClient

	authController := controllers.NewAuthController(grpcUserClient, grpcSessionClient, opts.PrivateKey, opts.PublicKey)
	mangaController := controllers.NewMangaController(grpcMangaClient, grpcChapterClient)
	
	//2. Middleware for public routes can be added here (if needed)

	//3. Route definition
	// AUTH ROUTES
	rg.POST("/login", authController.LoginByUsername)
	rg.POST("/signup", authController.SignupByUsername)
	rg.POST("/auth/refresh", authController.RefreshToken)

	// MANGA ROUTES
	rg.GET("/mangas", mangaController.ListMangas)
	rg.GET("/mangas/:id", mangaController.GetMangaDetail)
	rg.GET("/mangas/:id/chapters", mangaController.GetMangaChapters)

	// CHAPTER ROUTES
	rg.GET("/chapters/:chapter_id", mangaController.ReadChapter)
}
