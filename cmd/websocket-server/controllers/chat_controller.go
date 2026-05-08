package controllers

import (
	"log"
	"mangahub/internal/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatController struct {
	chatService websocket_services.ChatService
}

func NewChatController(chatService websocket_services.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
	}
}

func (h *ChatController) HandleWSChatTunnel(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Adjust for production
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	userIDVal, _ := c.Get("user_id")
	userID, ok := userIDVal.(string)
	if !ok {
		conn.WriteJSON(gin.H{"error": "invalid user context"})
		conn.Close()
		return
	}

	roomID := c.Query("manga_id")
	if userID == "" || roomID == "" {
		conn.WriteJSON(gin.H{"error": "user_id and manga_id are required"})
		conn.Close()
		return
	}

	h.chatService.HandleWSChatTunnel(conn, userID, roomID)
}
