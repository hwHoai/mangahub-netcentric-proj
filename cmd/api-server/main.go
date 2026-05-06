package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	routes "mangahub/cmd/api-server/routes"
	tcp_services_impl "mangahub/internal/tcp/impl"
	udp_services_impl "mangahub/internal/udp/impl"
	websocket_impl "mangahub/internal/websocket/impl"
	jwt_impl "mangahub/pkg/utils/jwt/impl"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load file .env	
	if err := godotenv.Load("../../.env"); err != nil {
		// log.Fatalf will stop the program
		log.Println("Warning: No .env file found, using environment variables if set")
	}
	tcpHost := os.Getenv("TCP_SERVER_HOST")
	tcpPort := os.Getenv("TCP_SERVER_PORT")
	if tcpHost == "" {
		tcpHost = "localhost"
	}
	if tcpPort == "" {
		tcpPort = "8082"
	}
	tcpAddr := fmt.Sprintf("%s:%s", tcpHost, tcpPort)

	udpHost := os.Getenv("UDP_SERVER_HOST")
	udpPort := os.Getenv("UDP_SERVER_PORT")
	if udpHost == "" {
		udpHost = "127.0.0.1"
	}
	if udpPort == "" {
		udpPort = "8083"
	}
	udpAddr := fmt.Sprintf("%s:%s", udpHost, udpPort)

	//2. Setup Router
	r := gin.Default()

	// 3. Middleware (CORS, Logger, Recovery...)
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// 4. Generate JWT key pair once at startup and keep private key in memory only for auth.
	jwtUtil := jwt_impl.NewJWTUtil(nil)
	privateKey, publicKey, err := jwtUtil.CreateRSAKeyPair(2048)
	if err != nil {
		log.Fatalf("failed to create JWT key pair: %v", err)
	}

	// 4.1 Broadcast public key to TCP server
	publicKeyPEM, _ := jwtUtil.StringifyPublicKeyPEM(publicKey)
	tcpKeySyncService := tcp_services_impl.NewTCPKeySyncService(tcpAddr)
	if err := tcpKeySyncService.SyncPublicKey(publicKeyPEM); err != nil {
		log.Printf("Warning: failed to broadcast public key to TCP server: %v", err)
	}

	// 4.2 Broadcast public key to UDP server
	udpKeySyncService := udp_services_impl.NewUDPKeySyncService(udpAddr)
	if err := udpKeySyncService.SyncPublicKey(publicKeyPEM); err != nil {
		log.Printf("Warning: failed to broadcast public key to UDP server: %v", err)
	}

	// 4.3 Broadcast public key to WebSocket server
	wsPort := os.Getenv("WS_SERVER_PORT")
	if wsPort == "" {
		wsPort = "8085"
	}
	wsAddr := "localhost:" + wsPort
	wsKeySyncService := websocket_impl.NewWSKeySyncService(wsAddr)
	if err := wsKeySyncService.SyncPublicKey(publicKeyPEM); err != nil {
		log.Printf("Warning: failed to broadcast public key to WebSocket server: %v", err)
	}
	
	// 5. Routes definition
	routes.SetupRoutes(r, privateKey, publicKey)
	privateKey = nil

	// 6. Configure HTTP Server
	port, srv := getServerConfiguration(r)
	
	// 7. Start Server in a goroutine
	go startServer(port, srv)

	// 8. Graceful Shutdown
	shutdownServer(srv)
}

func getServerConfiguration(r *gin.Engine) (string, *http.Server) {
	// 1. Setup HTTP Port
	port := ":" + os.Getenv("API_SERVER_PORT")
	if port == ":" {
		port = ":8081"
	}

	// 2. Create HTTP Server
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	return port, srv
}

func startServer(port string, srv *http.Server) {
	log.Printf("Server starting on port %s", port)
	
	// 2. Start Server
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Listen error: %s\n", err)
	}
}

func shutdownServer(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}