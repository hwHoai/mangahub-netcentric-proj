package udp_pool_impl

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"

	udp_pools "mangahub/cmd/udp-server/utils/pools"
	"mangahub/proto/user_manga"
)

type messageNotificationPool struct {
	mu         sync.RWMutex
	clients    map[string]*net.UDPAddr
	grpcClient user_manga.GRPCUserMangaServiceClient
}

var _ udp_pools.UDPPool = (*messageNotificationPool)(nil)

func NewMessageNotificationPool(grpcClient user_manga.GRPCUserMangaServiceClient) udp_pools.UDPPool {
	return &messageNotificationPool{
		clients:    make(map[string]*net.UDPAddr),
		grpcClient: grpcClient,
	}
}

func (p *messageNotificationPool) Register(userID string, addr *net.UDPAddr) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clients[userID] = addr
	log.Printf("UDP chat client registered for user %s at %s", userID, addr.String())
}

func (p *messageNotificationPool) Broadcast(conn *net.UDPConn, roomID string, payload map[string]interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// For chat, we also notify followers of the manga (roomID is mangaID)
	res, err := p.grpcClient.GetFollowers(ctx, &user_manga.GetFollowersRequest{
		MangaId: roomID,
	})
	if err != nil {
		log.Printf("Failed to get followers for chat room %s: %v", roomID, err)
		return
	}

	dataBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal UDP chat payload: %v", err)
		return
	}
	dataBytes = append(dataBytes, '\n')

	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, userID := range res.UserIds {
		if addr, exists := p.clients[userID]; exists {
			_, err := conn.WriteToUDP(dataBytes, addr)
			if err != nil {
				log.Printf("Failed to send UDP chat message to %s: %v", addr.String(), err)
			} else {
				log.Printf("Sent new message notification to user %s at %s", userID, addr.String())
			}
		}
	}
}
