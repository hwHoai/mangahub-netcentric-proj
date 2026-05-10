package main

import (
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
	"mangahub/pkg/logger"
)

func main() {
	// 1. Load env
	if err := godotenv.Load("../../.env"); err != nil {
		logger.Warn("No .env file found, using environment variables if set")
	}
	port := ":" + os.Getenv("GRPC_SERVER_PORT")
	if port == ":" {
		port = ":8084"
	}
	logger.Init(os.Getenv("ENV") == "prod", 0)
	logger.Info("gRPC Server starting...", "pid", os.Getpid())
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../../data/mangahub.db"
	}

	// 1. Database Connection
	database := &dbImpl.SqliteConnImpl{}
	dbConn, err := database.InitDB(dbPath)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
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
		logger.Error("Start gRPC Server error", "error", err)
		os.Exit(1)
	}

	// 6. Run Server
	logger.Info("gRPC Server is running", "port", port)
	if err := grpcServer.Serve(listener); err != nil {
		logger.Error("Start gRPC Server error", "error", err)
		os.Exit(1)
	}
}

// seedData retrieves manga data from MangaDex API and stores it in the database
// It fetches 50 manga per batch and processes 2 batches (100 manga total)
func seedData(db *gorm.DB) {
	logger.Info("=== Starting Manga Data Seeding ===")

	// Check if database already has manga data
	var count int64
	if err := db.Model(&struct{}{}).Table("mangas").Count(&count).Error; err == nil && count > 0 {
		logger.Info("Database already contains manga. Skipping seeding.", "count", count)
		return
	}

	// Initialize manga seeder
	mangaSeeder := seeder.NewMangaSeeder(db)

	// Seed manga data: 50 per batch, 2 batches = 100 manga
	// Adjust these values as needed
	if err := mangaSeeder.SeedMangaData(50, 2); err != nil {
		logger.Error("Error during manga seeding", "error", err)
		// Don't fatally fail if seeding fails - the server can still run
	}

	logger.Info("=== Manga Data Seeding Completed ===")
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
