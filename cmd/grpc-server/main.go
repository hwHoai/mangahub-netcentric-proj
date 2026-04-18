package main

import (
	"log"
	"mangahub/internal/grpc/services"
	"mangahub/proto/sample"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// 1. Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using environment variables if set")
	}
	port := ":" + os.Getenv("GRPC_SERVER_PORT")
	if port == ":" {
		port = ":8084"
	}
	// 2. Initialize gRPC Server
	grpcServer := grpc.NewServer()

	// 3. Register services
	sampleService := &services.GRPCSampleService{}
	sample.RegisterSampleServiceServer(grpcServer, sampleService)

	// 4. Listen for connections
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Start gRPC Server error: %v", err)
	}

	// 5. Run Server
	log.Printf("gRPC Server is running on port %s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Start gRPC Server error: %v", err)
	}
}