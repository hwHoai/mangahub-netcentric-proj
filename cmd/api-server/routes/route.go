package routes

import (
	"crypto/rsa"
	"log"
	manga_services_impl "mangahub/internal/manga/impl"
	udp_services_impl "mangahub/internal/udp/impl"
	user_services_impl "mangahub/internal/user/impl"
	"mangahub/pkg/clients"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) {
	// 0. Setup shared clients
	mangaDexClient := clients.NewMangaDexClient()
	udpNotificationServices, err := udp_services_impl.NewNotificationServicesImpl()
	if err != nil {
		log.Printf("Warning: failed to initialize UDP notification services: %v", err)
	}

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

	// 2. Initialize Services
	mangaService := manga_services_impl.NewMangaService(grpcMangaClient)
	chapterService := manga_services_impl.NewChapterService(grpcChapterClient, udpNotificationServices, mangaDexClient)
	userMangaService := user_services_impl.NewUserMangaService(grpcUserMangaClient)
	userService := user_services_impl.NewUserService(grpcUserClient)

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
		GRPCMessageClient: grpcMessageClient,
		PrivateKey:        privateKey,
		PublicKey:         publicKey,
	}
	SetupPublicRoutes(v1, public_route_opts)
 
	private_route_opts := &PrivateRouteOpts{
		PublicKey:            publicKey,
		UserMangaService:     userMangaService,
		UserService:            userService,
		GRPCUserClient:       grpcUserClient,
		GRPCSessionClient:    grpcSessionClient,
		ChapterService:       chapterService,
		TCPChapterSyncClient: tcpChapterSyncClient,
	}
	SetupPrivateRoutes(v1, private_route_opts)
}