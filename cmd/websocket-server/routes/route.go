package routes

import (
	"mangahub/cmd/websocket-server/controllers"
	"mangahub/cmd/websocket-server/middleware"
	ws_utils_pool_impl "mangahub/cmd/websocket-server/utils/pool/impl"
	"mangahub/internal/websocket/impl"
	"mangahub/pkg/clients"
	"mangahub/pkg/logger"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// 1. Initialize gRPC clients
	messageClient, _, err := clients.NewMessageGRPCClient()
	if err != nil {
		logger.Error("Failed to initialize message gRPC client", "error", err)
	}

	mangaClient, _, err := clients.NewMangaGRPCClient()
	if err != nil {
		logger.Error("Failed to initialize manga gRPC client", "error", err)
	}

	// 2. Initialize UDP client for notifications
	udpClient, err := clients.NewUDPNotificationClient()
	if err != nil {
		logger.Error("Failed to initialize UDP client", "error", err)
	}

	// 3. Initialize WebSocket Pool
	pool := ws_utils_pool_impl.NewChatPool()

	// 4. Init Services
	chatService := websocket_impl.NewWSChatService(pool, messageClient, mangaClient, udpClient)

	// 6. Init Controllers
	chatController := controllers.NewChatController(chatService)
	keySyncController := controllers.NewKeySyncController()

	// Internal route for public key sync
	r.POST("/impl/sync-public-key", keySyncController.SyncPublicKeyHandler)

	// WebSocket route with Auth Middleware
	wsGroup := r.Group("/ws")
	wsGroup.Use(middleware.AuthMiddleware())
	{
		wsGroup.GET("", chatController.HandleWSChatTunnel)
	}
}
