package main

import (
	"log"
	"os"

	ws_utils_pool_impl "mangahub/cmd/websocket-server/utils/pool/impl"
	udp_services_impl "mangahub/internal/udp/impl"
	"mangahub/cmd/websocket-server/handler"
	"mangahub/internal/websocket/impl"
	"mangahub/cmd/websocket-server/middleware"
	"mangahub/pkg/clients"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Warning: No .env file found, using environment variables if set")
	}

	port := os.Getenv("WS_SERVER_PORT")
	if port == "" {
		port = "8085"
	}

	// 2. Initialize gRPC clients
	messageClient, grpcMsgConn, err := clients.NewMessageGRPCClient()
	if err != nil {
		log.Fatalf("Failed to initialize message gRPC client: %v", err)
	}
	defer grpcMsgConn.Close()

	mangaClient, grpcMangaConn, err := clients.NewMangaGRPCClient()
	if err != nil {
		log.Fatalf("Failed to initialize manga gRPC client: %v", err)
	}
	defer grpcMangaConn.Close()

	// 3. Initialize UDP client for notifications
	udpClient, err := udp_services_impl.NewNotificationServicesImpl()
	if err != nil {
		log.Printf("Warning: Failed to initialize UDP client: %v", err)
	} else {
		defer udpClient.Close()
	}

	// 4. Initialize WebSocket Pool & Service
	pool := ws_utils_pool_impl.NewChatPool()
	chatService := websocket_impl.NewWSChatService(pool, messageClient, mangaClient, udpClient)
	chatHandler := handler.NewChatHandler(chatService)
	keySyncHandler := handler.NewKeySyncHandler()

	// 5. Setup Gin
	r := gin.Default()
	
	// Internal route for public key sync
	r.POST("/impl/sync-public-key", keySyncHandler.SyncPublicKeyHandler)

	// WebSocket route with Auth Middleware
	wsGroup := r.Group("/ws")
	wsGroup.Use(middleware.AuthMiddleware())
	{
		wsGroup.GET("", chatHandler.HandleWSChatTunnel)
	}

	// 6. Run Server
	log.Printf("WebSocket Server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start WebSocket server: %v", err)
	}
}
