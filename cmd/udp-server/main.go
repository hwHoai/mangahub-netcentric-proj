package main

import (
	"log"
	"mangahub/cmd/udp-server/dispatch"
	"mangahub/cmd/udp-server/handler"
	pool_impl "mangahub/cmd/udp-server/utils/pool/impl"
	"mangahub/pkg/clients"
	"mangahub/pkg/logger"
	"mangahub/pkg/types"
	"net"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	//1. Load env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Warning: No .env file found, using environment variables if set")
	}
	port := ":" + os.Getenv("UDP_SERVER_PORT")
	if port == ":" {
		port = ":8083"
	}
	logger.Init(os.Getenv("ENV") == "prod", 0)
	logger.Info("UDP Server starting...", "pid", os.Getpid())

	//2. Setup gRPC client
	grpcUserMangaClient, grpcConn, err := clients.NewUserMangaGRPCClient()
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer grpcConn.Close()

	//3. Setup Dispatcher
	udpServer := dispatch.NewUDPServer()

	//4. Setup Pool and Handlers
	chapterPool := pool_impl.NewChapterNotificationPool(grpcUserMangaClient)
	messagePool := pool_impl.NewMessageNotificationPool(grpcUserMangaClient)
	notificationHandler := handler.NewNotificationHandler(chapterPool, messagePool)
	keySyncHandler := handler.NewKeySyncHandler()

	// Register handlers
	udpServer.RegisterHandler("chapter:req_client_register", notificationHandler.ClientRegisterHandler)
	udpServer.RegisterHandler("chapter:broadcast_chapter", notificationHandler.BroadcastChapterHandler)
	udpServer.RegisterHandler("chat:broadcast_message", notificationHandler.BroadcastMessageHandler)
	udpServer.RegisterHandler("chapter:ack_notification", notificationHandler.NotificationAckHandler)
	udpServer.RegisterHandler("pub_key:impl_sync_public_key", keySyncHandler.SyncPublicKeyHandler)

	// Benchmark handler (Ping-Pong)
	udpServer.RegisterHandler("benchmark:test_ping", func(s *dispatch.UDPServer, addr *net.UDPAddr, msg types.UDPMessage) {
		log.Printf("UDP Ping received from %v (ID: %s)", addr, string(msg.Payload))
	})

	//5. Resolve UDP address and Start
	udpServer.Start(port)

}