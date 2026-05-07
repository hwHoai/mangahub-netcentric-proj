package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"mangahub/cmd/tcp-server/dispatch"
	"mangahub/cmd/tcp-server/handler"
	"mangahub/cmd/tcp-server/utils/pool"
	pool_impl "mangahub/cmd/tcp-server/utils/pool/impl"
	"mangahub/pkg/clients"
	"mangahub/pkg/logger"
	"mangahub/pkg/types"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Warning: No .env file found, using environment variables if set")
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
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer grpcConn.Close()

	//3. Setup connection pool
	chapterSyncPool := pool_impl.NewChapterSyncPool(grpcUserMangaClient)

	//4. Init handler
	chapterSyncHandler := handler.NewChapterSyncHandler(chapterSyncPool)
	keySyncHandler := handler.NewKeySyncHandler()

	//5. Setup dispatcher
	dispatcher := dispatch.NewDispatcher()

	//6. Register handler
	// Chapter sync handler
	dispatcher.RegisterHandler("chapter_sync:req_register_client", chapterSyncHandler.RegisterConnectionHandler)
	dispatcher.RegisterHandler("chapter_sync:impl_broadcast_read", chapterSyncHandler.BroadcastReadHandler)

	// Key sync handler
	dispatcher.RegisterHandler("pub_key:impl_sync_public_key", keySyncHandler.SyncPublicKeyHandler)

	// Benchmark handler (Ping-Pong)
	dispatcher.RegisterHandler("benchmark:test_ping", func(conn net.Conn, payload any) {
		response := types.TCPMessage{
			Action:  "benchmark:res_pong",
			Payload: json.RawMessage(`{"status": "ok", "msg": "PONG"}`),
		}
		data, _ := json.Marshal(response)
		fmt.Fprintln(conn, string(data))
	})

	dispatcher.RegisterHandler("benchmark:res_pong", func(conn net.Conn, payload any) {
		// Just a placeholder to show it's received on server logs
		// No action needed
	})

	//7. Open TCP port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Start TCP Server error: %v", err)
	}
	defer listener.Close()
	log.Printf("TCP server started on port %s", port)

	// 4. Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// 5. Handle connection in a new goroutine to allow multiple clients
		go handleTCPConnection(conn, dispatcher, []pools.ConnectionPool{
			chapterSyncPool,
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
		log.Printf("Connection closed: %s", conn.RemoteAddr())
	}()

	conn.SetReadDeadline(time.Now().Add(tcpIdleTimeout))
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		conn.SetReadDeadline(time.Now().Add(tcpIdleTimeout))

		var msg types.TCPMessage
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			log.Printf("Error decoding message from %s: %v (raw: %s)", conn.RemoteAddr(), err, scanner.Text())
			continue
		}
		log.Printf("Received action: %s", msg.Action) 

		dispatcher.Dispatch(conn, msg)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Connection error from %s: %v", conn.RemoteAddr(), err)
	}
}