package routes

import (
	"crypto/rsa"
	"mangahub/cmd/api-server/controllers"
	"mangahub/cmd/api-server/middleware"
	"mangahub/proto/user_manga"

	"github.com/gin-gonic/gin"
)

type PrivateRouteOpts struct {
	PublicKey           *rsa.PublicKey
	GRPCUserMangaClient user_manga.GRPCUserMangaServiceClient
}

func SetupPrivateRoutes(rg *gin.RouterGroup, opts *PrivateRouteOpts) {
	authMiddleware := middleware.NewAuthMiddleware(opts.PublicKey)
	rg.Use(authMiddleware.Handler())

	// Initialize controller
	userMangaController := controllers.NewUserMangaController(opts.GRPCUserMangaClient)

	// USER MANGA ROUTES
	userMangas := rg.Group("/user/mangas")
	{
		userMangas.GET("/following", userMangaController.GetFollowingMangas)
		userMangas.POST("/:id/follow", userMangaController.FollowManga)
		userMangas.DELETE("/:id/follow", userMangaController.UnfollowManga)
	}

	// USER CHAPTER ROUTES
	userChapters := rg.Group("/user/chapters")
	{
		userChapters.POST("/:chapter_id/read", userMangaController.StoreReadingProgress)
	}
}