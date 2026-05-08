package routes

import (
	"crypto/rsa"
	"mangahub/cmd/api-server/controllers"
	"mangahub/internal/manga"
	"mangahub/internal/scrape"
	user_internal "mangahub/internal/user"
	"mangahub/proto/session"
	"mangahub/proto/user"
	"mangahub/proto/message"

	"github.com/gin-gonic/gin"
)

type PublicRouteOpts struct {
	GRPCUserClient    user.GRPCUserServiceClient
	GRPCSessionClient session.GRPCSessionServiceClient
	MangaService      manga.MangaService
	ChapterService    manga.ChapterService
	UserService       user_internal.UserService
	ScrapeService     scrape.ScrapeService
	GRPCMessageClient message.GRPCMessageServiceClient
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
	scrapeController := controllers.NewScrapeController(opts.ScrapeService)
	chatController := controllers.NewChatController(opts.GRPCMessageClient)
	
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
	rg.GET("/mangas/:id/messages", chatController.GetChatHistory)

	// CHAPTER ROUTES
	rg.GET("/chapters/:chapter_id", mangaController.ReadChapter)

	// EDUCATIONAL SCRAPING ROUTE
	rg.GET("/quotes", scrapeController.ScrapeQuotes)
}
