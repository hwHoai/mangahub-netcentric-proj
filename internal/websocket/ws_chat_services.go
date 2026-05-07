package websocket_services

import (
	"github.com/gorilla/websocket"
)

type ChatService interface {
	HandleWSChatTunnel(conn *websocket.Conn, userID string, roomID string)
}