package udp_pool_impl

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	benchmarks_prometheus "mangahub/benchmarks/prometheus"
	udp_pools "mangahub/cmd/udp-server/utils/pool"
	"mangahub/pkg/logger"
	"mangahub/proto/user_manga"
)

type chapterNotificationPool struct {
	mu                sync.RWMutex
	clients           map[string]*net.UDPAddr
	failCount         map[string]int
	lastSeen          map[string]time.Time
	pendingAcks       map[string]chan bool // key: notificationID:userID
	grpcClient        user_manga.GRPCUserMangaServiceClient
	prometheusMetrics *benchmarks_prometheus.Metrics
}

var _ udp_pools.UDPPool = (*chapterNotificationPool)(nil)

func NewChapterNotificationPool(grpcClient user_manga.GRPCUserMangaServiceClient, prometheusMetrics *benchmarks_prometheus.Metrics) udp_pools.UDPPool {
	p := &chapterNotificationPool{
		clients:           make(map[string]*net.UDPAddr),
		failCount:         make(map[string]int),
		lastSeen:          make(map[string]time.Time),
		pendingAcks:       make(map[string]chan bool),
		grpcClient:        grpcClient,
		prometheusMetrics: prometheusMetrics,
	}

	go p.startCleanupLoop()

	return p
}

func (p *chapterNotificationPool) startCleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		var stale []string
		p.mu.Lock()
		for userID, lastSeen := range p.lastSeen {
			if time.Since(lastSeen) > 5*time.Minute {
				stale = append(stale, userID)
			}
		}
		p.mu.Unlock()

		for _, userID := range stale {
			p.Unregister(userID, "TTL")
		}
	}
}

func (p *chapterNotificationPool) Register(userID string, addr *net.UDPAddr) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.clients[userID]; !exists {
		p.prometheusMetrics.ActiveConnections.Inc()
	}

	p.clients[userID] = addr
	p.lastSeen[userID] = time.Now()
	p.prometheusMetrics.ResponsesSent.Inc()
	p.prometheusMetrics.TotalRequests.Inc()
	p.failCount[userID] = 0
	logger.Info("UDP client registered", "userID", userID, "addr", addr.String())
}

func (p *chapterNotificationPool) Unregister(userID string, reason string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.clients[userID]; exists {
		delete(p.clients, userID)
		delete(p.failCount, userID)
		delete(p.lastSeen, userID)
		p.prometheusMetrics.ActiveConnections.Dec()
		logger.Info("UDP chapter client removed", "userID", userID, "reason", reason)
	}
}

func (p *chapterNotificationPool) Broadcast(conn *net.UDPConn, mangaID string, payload map[string]interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := p.grpcClient.GetFollowers(ctx, &user_manga.GetFollowersRequest{
		MangaId: mangaID,
	})
	if err != nil {
		logger.Error("Failed to get followers", "mangaID", mangaID, "error", err)
		return
	}

	// Generate a unique notification ID for ACK tracking if not present
	notificationID, ok := payload["id"].(string)
	if !ok || notificationID == "" {
		notificationID = fmt.Sprintf("notif_%d", time.Now().UnixNano())
		payload["id"] = notificationID
	}

	dataBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal UDP payload", "error", err)
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
		p.prometheusMetrics.ResponsesSent.Inc()
		p.prometheusMetrics.TotalRequests.Inc()
		if exists {
			go p.sendWithRetry(conn, userID, addr, notificationID, dataBytes)
		}
	}
}

func (p *chapterNotificationPool) sendWithRetry(conn *net.UDPConn, userID string, addr *net.UDPAddr, notificationID string, data []byte) {
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
			logger.Error("Failed to send UDP message", "userID", userID, "addr", addr.String(), "error", err)
			p.incrementFailCount(userID)
			return
		}

		logger.Info("Sent UDP notification", "userID", userID, "notificationID", notificationID, "attempt", attempt)

		select {
		case <-ackChan:
			logger.Info("Received ACK for UDP notification", "userID", userID, "notificationID", notificationID)
			p.resetFailCount(userID)

			// Update lastSeen on ACK
			p.mu.Lock()
			if _, exists := p.clients[userID]; exists {
				p.lastSeen[userID] = time.Now()
			}
			p.mu.Unlock()

			return
		case <-time.After(2 * time.Second):
			logger.Warn("UDP notification timeout, retrying...", "userID", userID, "notificationID", notificationID, "attempt", attempt)
		}

		p.prometheusMetrics.ResponsesSent.Inc()
		p.prometheusMetrics.TotalRequests.Inc()
	}

	logger.Error("UDP notification failed after 3 attempts", "userID", userID, "notificationID", notificationID)
	p.incrementFailCount(userID)
}

func (p *chapterNotificationPool) ProcessAck(notificationID string, userID string) {
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

func (p *chapterNotificationPool) incrementFailCount(userID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.failCount[userID]++
	if p.failCount[userID] >= 5 {
		if _, exists := p.clients[userID]; exists {
			delete(p.clients, userID)
			delete(p.failCount, userID)
			delete(p.lastSeen, userID)
			p.prometheusMetrics.ActiveConnections.Dec()
			logger.Info("UDP chapter client removed", "userID", userID, "reason", "max_failures")
		}
	}
}

func (p *chapterNotificationPool) resetFailCount(userID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.failCount[userID] = 0
}
