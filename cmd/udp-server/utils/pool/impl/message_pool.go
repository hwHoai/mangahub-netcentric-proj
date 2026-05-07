package udp_pool_impl

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	udp_pools "mangahub/cmd/udp-server/utils/pool"
	"mangahub/pkg/logger"
	"mangahub/proto/user_manga"
)

type messageNotificationPool struct {
	mu          sync.RWMutex
	clients     map[string]*net.UDPAddr
	failCount   map[string]int
	pendingAcks map[string]chan bool // key: notificationID:userID
	grpcClient  user_manga.GRPCUserMangaServiceClient
}

var _ udp_pools.UDPPool = (*messageNotificationPool)(nil)

func NewMessageNotificationPool(grpcClient user_manga.GRPCUserMangaServiceClient) udp_pools.UDPPool {
	return &messageNotificationPool{
		clients:     make(map[string]*net.UDPAddr),
		failCount:   make(map[string]int),
		pendingAcks: make(map[string]chan bool),
		grpcClient:  grpcClient,
	}
}

func (p *messageNotificationPool) Register(userID string, addr *net.UDPAddr) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clients[userID] = addr
	p.failCount[userID] = 0
	logger.Info("UDP chat client registered", "userID", userID, "addr", addr.String())
}

func (p *messageNotificationPool) Unregister(userID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.clients, userID)
	delete(p.failCount, userID)
	logger.Info("UDP chat client unregistered", "userID", userID)
}

func (p *messageNotificationPool) Broadcast(conn *net.UDPConn, roomID string, payload map[string]interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := p.grpcClient.GetFollowers(ctx, &user_manga.GetFollowersRequest{
		MangaId: roomID,
	})
	if err != nil {
		logger.Error("Failed to get followers for chat room", "roomID", roomID, "error", err)
		return
	}

	// Generate a unique notification ID for ACK tracking if not present
	notificationID, ok := payload["id"].(string)
	if !ok || notificationID == "" {
		notificationID = fmt.Sprintf("notif_chat_%d", time.Now().UnixNano())
		payload["id"] = notificationID
	}

	dataBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal UDP chat payload", "error", err)
		return
	}
	dataBytes = append(dataBytes, '\n')

	p.mu.RLock()
	userIds := make([]string, len(res.UserIds))
	copy(userIds, res.UserIds)
	p.mu.RUnlock()

	for _, userID := range userIds {
		p.mu.RLock()
		addr, exists := p.clients[userID]
		p.mu.RUnlock()

		if exists {
			go p.sendWithRetry(conn, userID, addr, notificationID, dataBytes)
		}
	}
}

func (p *messageNotificationPool) sendWithRetry(conn *net.UDPConn, userID string, addr *net.UDPAddr, notificationID string, data []byte) {
	ackKey := fmt.Sprintf("%s:%s", notificationID, userID)
	ackChan := make(chan bool, 1)

	p.mu.Lock()
	p.pendingAcks[ackKey] = ackChan
	p.mu.Unlock()

	defer func() {
		p.mu.Lock()
		delete(p.pendingAcks, ackKey)
		p.mu.Unlock()
	}()

	for attempt := 1; attempt <= 3; attempt++ {
		_, err := conn.WriteToUDP(data, addr)
		if err != nil {
			logger.Error("Failed to send UDP chat message", "userID", userID, "addr", addr.String(), "error", err)
			p.incrementFailCount(userID)
			return
		}

		logger.Info("Sent UDP chat notification", "userID", userID, "notificationID", notificationID, "attempt", attempt)

		select {
		case <-ackChan:
			logger.Info("Received ACK for UDP chat notification", "userID", userID, "notificationID", notificationID)
			p.resetFailCount(userID)
			return
		case <-time.After(2 * time.Second):
			logger.Warn("UDP chat notification timeout, retrying...", "userID", userID, "notificationID", notificationID, "attempt", attempt)
		}
	}

	logger.Error("UDP chat notification failed after 3 attempts", "userID", userID, "notificationID", notificationID)
	p.incrementFailCount(userID)
}

func (p *messageNotificationPool) ProcessAck(notificationID string, userID string) {
	ackKey := fmt.Sprintf("%s:%s", notificationID, userID)
	p.mu.RLock()
	ackChan, exists := p.pendingAcks[ackKey]
	p.mu.RUnlock()

	if exists {
		select {
		case ackChan <- true:
		default:
		}
	}
}

func (p *messageNotificationPool) incrementFailCount(userID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.failCount[userID]++
	if p.failCount[userID] >= 5 {
		logger.Warn("User exceeded max UDP chat failures, unregistering", "userID", userID)
		delete(p.clients, userID)
		delete(p.failCount, userID)
	}
}

func (p *messageNotificationPool) resetFailCount(userID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.failCount[userID] = 0
}
