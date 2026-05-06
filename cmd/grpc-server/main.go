package main

import (
	"log"

	grpc_services_impl "mangahub/internal/grpc/impl"
	dbImpl "mangahub/internal/database/impl"
	"mangahub/pkg/seeder"
	"mangahub/proto/chapter"
	"mangahub/proto/manga"
	"mangahub/proto/message"
	"mangahub/proto/session"
	"mangahub/proto/user"
	"mangahub/proto/user_manga"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/gorm"
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
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../../data/mangahub.db"
	}

	// 1. Database Connection
	database := &dbImpl.SqliteConnImpl{}
	dbConn, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 2. Seed manga data from MangaDex API
	// seedData(dbConn)

	// 3. Initialize gRPC Server
	grpcServer := grpc.NewServer()

	// 4. Register services
	registerServices(grpcServer, dbConn)

	// 5. Listen for connections
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Start gRPC Server error: %v", err)
	}

	// 6. Run Server
	log.Printf("gRPC Server is running on port %s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Start gRPC Server error: %v", err)
	}
}

// seedData retrieves manga data from MangaDex API and stores it in the database
// It fetches 50 manga per batch and processes 2 batches (100 manga total)
func seedData(db *gorm.DB) {
	log.Println("=== Starting Manga Data Seeding ===")

	// Check if database already has manga data
	var count int64
	if err := db.Model(&struct{}{}).Table("mangas").Count(&count).Error; err == nil && count > 0 {
		log.Printf("Database already contains %d manga. Skipping seeding.", count)
		return
	}

	// Initialize manga seeder
	mangaSeeder := seeder.NewMangaSeeder(db)

	// Seed manga data: 50 per batch, 2 batches = 100 manga
	// Adjust these values as needed
	if err := mangaSeeder.SeedMangaData(50, 2); err != nil {
		log.Printf("Error during manga seeding: %v", err)
		// Don't fatally fail if seeding fails - the server can still run
	}

	log.Println("=== Manga Data Seeding Completed ===")
}

// registerServices registers all gRPC services with the server.
func registerServices(grpcServer *grpc.Server, db *gorm.DB) {
	userService := grpc_services_impl.NewGRPCUserService(db)
	user.RegisterGRPCUserServiceServer(grpcServer, userService)

	sessionService := grpc_services_impl.NewGRPCSessionService(db)
	session.RegisterGRPCSessionServiceServer(grpcServer, sessionService)

	mangaService := grpc_services_impl.NewGRPCMangaService(db)
	manga.RegisterGRPCMangaServiceServer(grpcServer, mangaService)

	userMangaService := grpc_services_impl.NewGRPCUserMangaService(db)
	user_manga.RegisterGRPCUserMangaServiceServer(grpcServer, userMangaService)

	chapterService := grpc_services_impl.NewGRPCChapterService(db)
	chapter.RegisterGRPCChapterServiceServer(grpcServer, chapterService)

	messageService := grpc_services_impl.NewGRPCMessageService(db)
	message.RegisterGRPCMessageServiceServer(grpcServer, messageService)
}
