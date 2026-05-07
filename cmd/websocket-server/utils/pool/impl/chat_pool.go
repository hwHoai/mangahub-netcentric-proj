package ws_utils_pool_impl

import (
	"sync"

	ws_utils_pool "mangahub/cmd/websocket-server/utils/pool"
	"mangahub/pkg/logger"
)

type ChatPoolImpl struct {
	mu      sync.RWMutex
	clients map[string]map[string]*ws_utils_pool.Client // roomID -> userID -> Client
}

var _ ws_utils_pool.ChatPool = (*ChatPoolImpl)(nil)

func NewChatPool() ws_utils_pool.ChatPool {
	return &ChatPoolImpl{
		clients: make(map[string]map[string]*ws_utils_pool.Client),
	}
}

func (p *ChatPoolImpl) Register(client *ws_utils_pool.Client) {
	p.mu.Lock()
	if p.clients[client.RoomID] == nil {
		p.clients[client.RoomID] = make(map[string]*ws_utils_pool.Client)
	}
	p.clients[client.RoomID][client.UserID] = client
	p.mu.Unlock()

	logger.Info("User joined room", "userID", client.UserID, "roomID", client.RoomID)

	// Broadcast join message
	p.Broadcast(ws_utils_pool.Message{
		RoomID:  client.RoomID,
		Content: "User " + client.UserID + " joined the room",
		Sender:  "system",
	})
}

func (p *ChatPoolImpl) Unregister(client *ws_utils_pool.Client) {
	p.mu.Lock()
	if p.clients[client.RoomID] != nil {
		delete(p.clients[client.RoomID], client.UserID)
		if len(p.clients[client.RoomID]) == 0 {
			delete(p.clients, client.RoomID)
		}
	}
	p.mu.Unlock()

	logger.Info("User left room", "userID", client.UserID, "roomID", client.RoomID)

	// Broadcast leave message
	p.Broadcast(ws_utils_pool.Message{
		RoomID:  client.RoomID,
		Content: "User " + client.UserID + " left the room",
		Sender:  "system",
	})
}

func (p *ChatPoolImpl) Broadcast(msg ws_utils_pool.Message) {
	p.mu.RLock()
	clients, ok := p.clients[msg.RoomID]
	p.mu.RUnlock()

	if ok {
		for _, client := range clients {
			if err := client.Conn.WriteJSON(msg); err != nil {
				logger.Error("Error broadcasting message", "userID", client.UserID, "error", err)
			}
		}
	}
}
