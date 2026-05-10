package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"mangahub/cmd/tcp-server/dispatch"
	"mangahub/cmd/tcp-server/handler"
	pools "mangahub/cmd/tcp-server/utils/pool"
	pool_impl "mangahub/cmd/tcp-server/utils/pool/impl"
	"mangahub/pkg/clients"
	"mangahub/pkg/logger"
	"mangahub/pkg/types"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Prometheus Exporter
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Println("Prometheus metrics available at http://localhost:2112/metrics")
		http.ListenAndServe(":2112", nil)
	}()

	// 1. Load env
	if err := godotenv.Load("../../.env"); err != nil {
		logger.Warn("No .env file found, using environment variables if set")
	}
	port := ":" + os.Getenv("TCP_SERVER_PORT")
	if port == ":" {
		port = ":8082"
	}
	logger.Init(os.Getenv("ENV") == "prod", 0)
	logger.Info("TCP Server starting...", "pid", os.Getpid())

	//2. Setup gRPC client
	grpcUserMangaClient, grpcConn, err := clients.NewUserMangaGRPCClient()
	if err != nil {
		logger.Error("Failed to create gRPC client", "error", err)
		os.Exit(1)
	}
	defer grpcConn.Close()

	//3. Setup connection pool
	chapterSyncPool := pool_impl.NewChapterSyncPool(grpcUserMangaClient)
	benchmarkPool := pool_impl.NewBenchmarkPool()

	//4. Init handler
	chapterSyncHandler := handler.NewChapterSyncHandler(chapterSyncPool)
	keySyncHandler := handler.NewKeySyncHandler()
	benchmarkHandler := handler.NewBenchmarkHandler(benchmarkPool)

	//5. Setup dispatcher
	dispatcher := dispatch.NewDispatcher()

	//6. Register handler
	// Chapter sync handler
	dispatcher.RegisterHandler("chapter_sync:req_register_client", chapterSyncHandler.RegisterConnectionHandler)
	dispatcher.RegisterHandler("chapter_sync:impl_broadcast_read", chapterSyncHandler.BroadcastReadHandler)

	// Key sync handler
	dispatcher.RegisterHandler("pub_key:impl_sync_public_key", keySyncHandler.SyncPublicKeyHandler)

	// Benchmark handler (Ping-Pong)
	dispatcher.RegisterHandler("benchmark:test_register", benchmarkHandler.RegisterHandler)
	dispatcher.RegisterHandler("benchmark:test_ping", benchmarkHandler.PingHandler)
	// dispatcher.RegisterHandler("benchmark:res_pong", benchmarkHandler.PongHandler)

	//7. Open TCP port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("Start TCP Server error", "error", err)
		os.Exit(1)
	}
	defer listener.Close()
	logger.Info("TCP server started", "port", port)

	// 4. Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("Error accepting connection", "error", err)
			continue
		}
		go handleTCPConnection(conn, dispatcher, []pools.ConnectionPool{
			chapterSyncPool,
			benchmarkPool,
		})
	}
}

const tcpIdleTimeout = 5 * time.Minute

func handleTCPConnection(conn net.Conn, dispatcher *dispatch.Dispatcher, pools []pools.ConnectionPool) {
	defer func() {
		for _, pool := range pools {
			go pool.Unregister(conn)
		}
		conn.Close()
		logger.Info("Connection closed", "addr", conn.RemoteAddr().String())
	}()

	conn.SetReadDeadline(time.Now().Add(tcpIdleTimeout))
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		conn.SetReadDeadline(time.Now().Add(tcpIdleTimeout))

		var msg types.TCPMessage
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			logger.Error("Error decoding message", "addr", conn.RemoteAddr().String(), "error", err, "raw", scanner.Text())
			continue
		}
		// logger.Info("Received action", "action", msg.Action) // High frequency - disabled for performance
		dispatcher.Dispatch(conn, msg)
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Connection error", "addr", conn.RemoteAddr().String(), "error", err)
	}
}