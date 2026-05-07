package routes

import (
	"crypto/rsa"
	"mangahub/cmd/api-server/controllers"
	"mangahub/cmd/api-server/middleware"
	"mangahub/internal/manga"
	tcp_services "mangahub/internal/tcp"
	user_internal "mangahub/internal/user"
	"mangahub/proto/session"
	user_proto "mangahub/proto/user"

	"github.com/gin-gonic/gin"
)

type PrivateRouteOpts struct {
	PublicKey            *rsa.PublicKey
	PrivateKey           *rsa.PrivateKey
	UserMangaService     user_internal.UserMangaService
	GRPCUserClient       user_proto.GRPCUserServiceClient
	GRPCSessionClient    session.GRPCSessionServiceClient
	MangaService         manga.MangaService
	ChapterService       manga.ChapterService
	UserService          user_internal.UserService
	TCPChapterSyncClient tcp_services.TCPChapterSyncServices
}

func SetupPrivateRoutes(rg *gin.RouterGroup, opts *PrivateRouteOpts) {
	authMiddleware := middleware.NewAuthMiddleware(opts.PublicKey)
	
	// Create a private group that requires authentication
	private := rg.Group("/")
	private.Use(authMiddleware.Handler())

	// Initialize controller
	userMangaController := controllers.NewUserMangaController(
		opts.UserMangaService,
		opts.ChapterService,
		opts.TCPChapterSyncClient,
	)

	// USER MANGA ROUTES
	userMangas := private.Group("/user/mangas")
	{
		userMangas.GET("/following", userMangaController.GetFollowingMangas)
		userMangas.POST("/:id/follow", userMangaController.FollowManga)
		userMangas.DELETE("/:id/follow", userMangaController.UnfollowManga)
	}

	mangaController := controllers.NewMangaController(nil, opts.ChapterService)
	private.POST("/mangas/:id/chapters", mangaController.CreateNewChapter)

	private.GET("/user/history", userMangaController.GetReadingHistory)

	// USER CHAPTER ROUTES
	userChapters := private.Group("/user/chapters")
	{
		userChapters.GET("/:chapter_id", userMangaController.ReadChapterWithDevicesSync)
	}

	// AUTH PROFILE ROUTES
	authController := controllers.NewAuthController(opts.GRPCUserClient, opts.GRPCSessionClient, opts.UserService, nil, opts.PublicKey)
	private.GET("/auth/me", authController.GetMe)
}