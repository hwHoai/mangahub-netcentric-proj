package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	routes "mangahub/cmd/api-server/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load file .env trước khi làm bất cứ việc gì khác
	
	if err := godotenv.Load("../../.env"); err != nil {
		// log.Fatalf will stop the program
		log.Println("Warning: No .env file found, using environment variables if set")
	}

	//2. Setup Router
	r := gin.Default()

	// 3. Middleware (CORS, Logger, Recovery...)
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	
	// 4. Routes definition
	routes.SetupRoutes(r)

	// 5. Configure HTTP Server
	port, srv := getServerConfiguration(r)
	
	// 6. Start Server in a goroutine
	go startServer(port, srv)

	// 7. Graceful Shutdown (Bonus điểm cao)
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