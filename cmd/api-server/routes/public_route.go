package routes

import (
	"crypto/rsa"
	"mangahub/cmd/api-server/controllers"
	"mangahub/internal/manga"
	user_internal "mangahub/internal/user"
	"mangahub/proto/session"
	"mangahub/proto/user"

	"github.com/gin-gonic/gin"
)

type PublicRouteOpts struct {
	GRPCUserClient    user.GRPCUserServiceClient
	GRPCSessionClient session.GRPCSessionServiceClient
	MangaService      manga.MangaService
	ChapterService    manga.ChapterService
	UserService       user_internal.UserService
	PrivateKey        *rsa.PrivateKey
	PublicKey         *rsa.PublicKey
}

func SetupPublicRoutes(rg *gin.RouterGroup, opts *PublicRouteOpts) {
	//1. Handler definition (if needed)
	grpcUserClient := opts.GRPCUserClient
	grpcSessionClient := opts.GRPCSessionClient
	mangaService := opts.MangaService
	chapterService := opts.ChapterService

	authController := controllers.NewAuthController(grpcUserClient, grpcSessionClient, opts.UserService, opts.PrivateKey, opts.PublicKey)
	mangaController := controllers.NewMangaController(mangaService, chapterService)
	
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
