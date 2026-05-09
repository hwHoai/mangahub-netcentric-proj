package main

import (
	"mangahub/cmd/udp-server/dispatch"
	"mangahub/cmd/udp-server/handler"
	pool_impl "mangahub/cmd/udp-server/utils/pool/impl"
	"mangahub/pkg/clients"
	"mangahub/pkg/logger"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	//1. Load env
	if err := godotenv.Load("../../.env"); err != nil {
		logger.Warn("No .env file found, using environment variables if set")
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
		logger.Error("Failed to create gRPC client", "error", err)
		os.Exit(1)
	}
	defer grpcConn.Close()

	//3. Setup Dispatcher
	udpServer := dispatch.NewUDPServer()

	//4. Setup Pool
	chapterPool := pool_impl.NewChapterNotificationPool(grpcUserMangaClient)
	messagePool := pool_impl.NewMessageNotificationPool(grpcUserMangaClient)
	benchmarkPool := pool_impl.NewBenchmarkPool()

	//5. Setup Handlers
	notificationHandler := handler.NewNotificationHandler(chapterPool, messagePool)
	keySyncHandler := handler.NewKeySyncHandler()
	benchmarkHandler := handler.NewBenchmarkHandler(benchmarkPool)

	//6. Register handlers
	udpServer.RegisterHandler("chapter:req_client_register", notificationHandler.ClientRegisterHandler)
	udpServer.RegisterHandler("chapter:broadcast_chapter", notificationHandler.BroadcastChapterHandler)
	udpServer.RegisterHandler("chat:broadcast_message", notificationHandler.BroadcastMessageHandler)
	udpServer.RegisterHandler("chapter:ack_notification", notificationHandler.NotificationAckHandler)
	udpServer.RegisterHandler("pub_key:impl_sync_public_key", keySyncHandler.SyncPublicKeyHandler)

	// Benchmark handler (Ping-Pong & Broadcast)
	udpServer.RegisterHandler("benchmark:test_ping", benchmarkHandler.PingHandler)
	udpServer.RegisterHandler("benchmark:trigger_broadcast", benchmarkHandler.BroadcastHandler)

	//7. Resolve UDP address and Start
	udpServer.Start(port)

}