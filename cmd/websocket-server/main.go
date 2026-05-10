package main

import (
	"os"

	"mangahub/cmd/websocket-server/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"mangahub/pkg/logger"
)

func main() {
	// 1. Load env
	if err := godotenv.Load("../../.env"); err != nil {
		logger.Error("Failed to load .env file", "error", err)
	}

	port := os.Getenv("WS_SERVER_PORT")
	if port == "" {
		port = "8085"
	}
	logger.Init(os.Getenv("ENV") == "prod", 0)
	logger.Info("WebSocket Server starting...", "pid", os.Getpid())

	// 2. Setup Gin & Routes
	r := gin.Default()
	routes.SetupRoutes(r)

	// 3. Run Server
	logger.Info("WebSocket Server is running", "port", port)
	if err := r.Run(":" + port); err != nil {
		logger.Error("Failed to start WebSocket server", "error", err)
	}
}
