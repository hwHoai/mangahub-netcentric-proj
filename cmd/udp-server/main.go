package main

import (
	"log"
	"mangahub/cmd/udp-server/handler"
	"mangahub/cmd/udp-server/dispatch"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	//1. Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using environment variables if set")
	}
	port := ":" + os.Getenv("UDP_SERVER_PORT")
	if port == ":" {
		port = ":8083"
	}

	//2. Setup Dispatcher
	udpServer := dispatch.NewUDPServer()
	udpServer.RegisterHandler("example", handler.ExampleUDPHandler)

	//2. Resolve UDP address
	udpServer.Start(port)

}