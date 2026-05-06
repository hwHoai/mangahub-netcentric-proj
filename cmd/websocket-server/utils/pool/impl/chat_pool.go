package ws_utils_pool_impl

import (
	"log"
	"sync"

	ws_utils_pool "mangahub/cmd/websocket-server/utils/pool"
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
	defer p.mu.Unlock()

	if p.clients[client.RoomID] == nil {
		p.clients[client.RoomID] = make(map[string]*ws_utils_pool.Client)
	}
	p.clients[client.RoomID][client.UserID] = client
	log.Printf("User %s joined room %s", client.UserID, client.RoomID)
}

func (p *ChatPoolImpl) Unregister(client *ws_utils_pool.Client) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.clients[client.RoomID] != nil {
		delete(p.clients[client.RoomID], client.UserID)
		if len(p.clients[client.RoomID]) == 0 {
			delete(p.clients, client.RoomID)
		}
	}
	log.Printf("User %s left room %s", client.UserID, client.RoomID)
}

func (p *ChatPoolImpl) Broadcast(msg ws_utils_pool.Message) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if clients, ok := p.clients[msg.RoomID]; ok {
		for _, client := range clients {
			if err := client.Conn.WriteJSON(msg); err != nil {
				log.Printf("Error broadcasting to user %s: %v", client.UserID, err)
			}
		}
	}
}
