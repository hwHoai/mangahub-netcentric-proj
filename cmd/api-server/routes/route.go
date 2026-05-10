package routes

import (
	"crypto/rsa"
	manga_services_impl "mangahub/internal/manga/impl"
	scrape_services_impl "mangahub/internal/scrape/impl"
	user_services_impl "mangahub/internal/user/impl"
	"mangahub/pkg/clients"
	"mangahub/pkg/logger"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) {
	// 0. Setup shared clients
	mangaDexClient := clients.NewMangaDexClient()

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

	grpcMessageClient, _, err := clients.NewMessageGRPCClient()
	if err != nil {
		panic(err)
	}

	tcpChapterSyncClient := clients.NewTCPChapterSyncClient()

	//2. Init UDP services
	udpNotificationClient, err := clients.NewUDPNotificationClient()
	if err != nil {
		logger.Warn("failed to initialize UDP notification services", "error", err)
	}

	// 2. Initialize Services
	mangaService := manga_services_impl.NewMangaService(grpcMangaClient)
	chapterService := manga_services_impl.NewChapterService(grpcChapterClient, udpNotificationClient, mangaDexClient)
	userMangaService := user_services_impl.NewUserMangaService(grpcUserMangaClient)
	userService := user_services_impl.NewUserService(grpcUserClient)
	scrapeService := scrape_services_impl.NewScrapeService()

	//3. Route definition
	v1 := r.Group("api/v1")
	
	// Health Check
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	public_route_opts := &PublicRouteOpts{
		GRPCUserClient:    grpcUserClient,
		GRPCSessionClient: grpcSessionClient,
		MangaService:      mangaService,
		ChapterService:    chapterService,
		UserService:       userService,
		ScrapeService:     scrapeService,
		GRPCMessageClient: grpcMessageClient,
		PrivateKey:        privateKey,
		PublicKey:         publicKey,
	}
	SetupPublicRoutes(v1, public_route_opts)
 
	private_route_opts := &PrivateRouteOpts{
		PublicKey:            publicKey,
		UserMangaService:     userMangaService,
		UserService:          userService,
		GRPCUserClient:       grpcUserClient,
		GRPCSessionClient:    grpcSessionClient,
		ChapterService:       chapterService,
		TCPChapterSyncClient: tcpChapterSyncClient,
	}
	SetupPrivateRoutes(v1, private_route_opts)
}