package routes

import (
	"crypto/rsa"
	"mangahub/cmd/api-server/controllers"
	"mangahub/cmd/api-server/middleware"
	tcp_services "mangahub/internal/tcp"
	"mangahub/proto/chapter"
	"mangahub/proto/session"
	"mangahub/proto/user"
	"mangahub/proto/user_manga"

	"github.com/gin-gonic/gin"
)

type PrivateRouteOpts struct {
	PublicKey            *rsa.PublicKey
	GRPCUserMangaClient  user_manga.GRPCUserMangaServiceClient
	GRPCUserClient       user.GRPCUserServiceClient
	GRPCSessionClient    session.GRPCSessionServiceClient
	GRPCChapterClient    chapter.GRPCChapterServiceClient
	TCPChapterSyncClient tcp_services.TCPChapterSyncServices
}

func SetupPrivateRoutes(rg *gin.RouterGroup, opts *PrivateRouteOpts) {
	authMiddleware := middleware.NewAuthMiddleware(opts.PublicKey)
	
	// Create a private group that requires authentication
	private := rg.Group("/")
	private.Use(authMiddleware.Handler())

	// Initialize controller
	userMangaController := controllers.NewUserMangaController(
		opts.GRPCUserMangaClient,
		opts.GRPCChapterClient,
		opts.TCPChapterSyncClient,
	)

	// USER MANGA ROUTES
	userMangas := private.Group("/user/mangas")
	{
		userMangas.GET("/following", userMangaController.GetFollowingMangas)
		userMangas.POST("/:id/follow", userMangaController.FollowManga)
		userMangas.DELETE("/:id/follow", userMangaController.UnfollowManga)
	}

	private.GET("/user/history", userMangaController.GetReadingHistory)

	// USER CHAPTER ROUTES
	userChapters := private.Group("/user/chapters")
	{
		userChapters.GET("/:chapter_id", userMangaController.ReadChapter)
	}

	// AUTH PROFILE ROUTES
	authController := controllers.NewAuthController(opts.GRPCUserClient, opts.GRPCSessionClient, nil, opts.PublicKey)
	private.GET("/auth/me", authController.GetMe)
}