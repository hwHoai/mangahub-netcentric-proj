package main

import (
	"log"

	"mangahub/internal/grpc/impl"
	dbImpl "mangahub/pkg/database/impl"
	"mangahub/proto/session"
	"mangahub/proto/user"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// 1. Load env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Warning: No .env file found, using environment variables if set")
	}
	port := ":" + os.Getenv("GRPC_SERVER_PORT")
	if port == ":" {
		port = ":8084"
	}

	// 1. Database Connection
	database := &dbImpl.SqliteConnImpl{}
	dbConn, err := database.InitDB(os.Getenv("DB_PATH"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 2. Initialize gRPC Server
	grpcServer := grpc.NewServer()

	// 3. Register services
	userService := &grpc_services_impl.GRPCUserService{DBConn: dbConn}
	user.RegisterGRPCUserServiceServer(grpcServer, userService)

	sessionService := &grpc_services_impl.GRPCSessionService{DBConn: dbConn}
	session.RegisterGRPCSessionServiceServer(grpcServer, sessionService)

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
