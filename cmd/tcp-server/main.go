package main

import (
	"bufio"
	"encoding/json"
	"log"
	"mangahub/cmd/tcp-server/dispatch"
	"mangahub/cmd/tcp-server/handler"
	"net"
	"os"

	"github.com/joho/godotenv"
)

// Cấu trúc gói tin chung
type Message struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}

func main() {
	// 1. Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using environment variables if set")
	}
	port := ":" + os.Getenv("TCP_SERVER_PORT")
	if port == ":" {
		port = ":8082"
	}

	//2. Setup Dispatcher
	dispatcher := dispatch.NewDispatcher()
	dispatcher.RegisterHandler("example", handler.ExampleTCPHandler)

	// 3. Open TCP port
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
		go handleTCPConnection(conn, dispatcher)
	}
}

func handleTCPConnection(conn net.Conn, dispatcher *dispatch.Dispatcher) {
	defer conn.Close()
	// Init decoder to read JSON messages from the connection
	scanner := bufio.NewScanner(conn)

	// Infinite loop to read messages from the client
	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			log.Printf("Error decoding message: %v", err)
			return
		}
		log.Printf("Received action: %s with payload: %s", msg.Action, string(msg.Payload))

		// Dispatch message to the appropriate handler based on the action
		dispatcher.Dispatch(conn, msg.Action, msg.Payload)
	}
}