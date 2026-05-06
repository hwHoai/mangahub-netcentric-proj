package handler

import (
	"log"
	"mangahub/internal/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	chatService websocket_services.ChatService
}

func NewChatHandler(chatService websocket_services.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (h *ChatHandler) HandleWSChatTunnel(c *gin.Context) {
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

	userIDVal, _ := c.Get("userID")
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
