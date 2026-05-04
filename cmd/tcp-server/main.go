package main

import (
	"bufio"
	"encoding/json"
	"log"
	"mangahub/cmd/tcp-server/dispatch"
	"mangahub/cmd/tcp-server/handler"
	"mangahub/cmd/tcp-server/utils/pools"
	pool_impl "mangahub/cmd/tcp-server/utils/pools/impl"
	"mangahub/pkg/clients"
	"mangahub/pkg/types"
	"net"
	"os"

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

func handleTCPConnection(conn net.Conn, dispatcher *dispatch.Dispatcher, pools []pools.ConnectionPool) {
	defer func() {
		for _, pool := range pools {
			go pool.Unregister(conn)
		}
		conn.Close()
	}()
	// Init decoder to read JSON messages from the connection
	scanner := bufio.NewScanner(conn)

	// Infinite loop to read messages from the client
	for scanner.Scan() {
		var msg types.TCPMessage
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			log.Printf("Error decoding message: %v", err)
			return
		}
		log.Printf("Received action: %s with payload: %s", msg.Action, string(msg.Payload))

		// Dispatch message to the appropriate handler based on the action
		dispatcher.Dispatch(conn, msg)
	}
}