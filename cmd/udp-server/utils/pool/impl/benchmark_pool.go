package udp_pool_impl

import (
	"crypto/sha256"
	"encoding/json"
	benchmarks_prometheus "mangahub/benchmarks/prometheus"
	"mangahub/pkg/logger"
	"mangahub/pkg/types"
	"net"
	"sync"
	"sync/atomic"
	"time"

	udp_pools "mangahub/cmd/udp-server/utils/pool"
)

type benchmarkPool struct {
	clients           map[string]*net.UDPAddr
	lastSeen          map[string]time.Time
	acks              sync.Map // Lưu trữ trạng thái ACK: key là "notifID:userID"
	mu                sync.RWMutex
	prometheusMetrics *benchmarks_prometheus.Metrics
}

var _ udp_pools.UDPPool = (*benchmarkPool)(nil)

func NewBenchmarkPool(prometheusMetrics *benchmarks_prometheus.Metrics) udp_pools.UDPPool {
	p := &benchmarkPool{
		clients:           make(map[string]*net.UDPAddr),
		lastSeen:          make(map[string]time.Time),
		prometheusMetrics: prometheusMetrics,
	}

	// Start cleanup goroutine
	go p.startCleanupLoop()

	return p
}

func (p *benchmarkPool) startCleanupLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var stale []string
		p.mu.Lock()
		for userID, lastSeen := range p.lastSeen {
			if time.Since(lastSeen) > 2*time.Minute {
				stale = append(stale, userID)
			}
		}
		p.mu.Unlock()

		for _, userID := range stale {
			p.Unregister(userID, "TTL")
		}
	}
}

func (p *benchmarkPool) Register(userID string, addr *net.UDPAddr) {
	// 1. Simulate JWT RSA256 Parsing & Validation (CPU Intensive)
	secret := []byte("rsa_public_key_simulation_" + userID)
	for i := 0; i < 1000; i++ {
		h := sha256.Sum256(secret)
		secret = h[:]
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.clients[userID]; !exists {
		p.prometheusMetrics.ActiveConnections.Inc()
	}

	p.clients[userID] = addr
	p.lastSeen[userID] = time.Now()
	p.prometheusMetrics.TotalRequests.Inc()
}

func (p *benchmarkPool) Unregister(userID string, reason string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.clients[userID]; exists {
		delete(p.clients, userID)
		delete(p.lastSeen, userID)
		p.prometheusMetrics.ActiveConnections.Dec()
		logger.Info("UDP benchmark client removed", "userID", userID, "reason", reason)
	}
}

func (p *benchmarkPool) Broadcast(conn *net.UDPConn, id string, payload map[string]interface{}) {
	p.mu.RLock()
	clientsToBroadcast := make(map[string]*net.UDPAddr, len(p.clients))
	for k, v := range p.clients {
		clientsToBroadcast[k] = v
	}
	p.mu.RUnlock()

	payloadBytes, _ := json.Marshal(payload)
	broadcastMsg := types.UDPMessage{
		Action:  "benchmark:test_broadcast",
		Payload: json.RawMessage(payloadBytes),
	}
	dataBytes, _ := json.Marshal(broadcastMsg)

	start := time.Now()
	successCount := int64(0)
	var wg sync.WaitGroup

	for userID, addr := range clientsToBroadcast {
		wg.Add(1)
		go func(uID string, target *net.UDPAddr) {
			defer wg.Done()

			ackKey := id + ":" + uID
			p.acks.Delete(ackKey) // Reset ACK state for this notification

			maxRetries := 3
			for retry := 0; retry < maxRetries; retry++ {
				// 1. Gửi gói tin
				conn.WriteToUDP(dataBytes, target)

				// 2. Đợi ACK (timeout 100ms cho mỗi lần thử)
				for wait := 0; wait < 10; wait++ {
					if _, received := p.acks.Load(ackKey); received {
						atomic.AddInt64(&successCount, 1)
						p.prometheusMetrics.ResponsesSent.Inc()
						p.prometheusMetrics.TotalRequests.Inc()

						// Update lastSeen on ACK
						p.mu.Lock()
						if _, exists := p.clients[uID]; exists {
							p.lastSeen[uID] = time.Now()
						}
						p.mu.Unlock()

						return // Đã nhận được ACK, thoát!
					}
					time.Sleep(10 * time.Millisecond)
				}

				// Nếu hết 100ms chưa có ACK -> Vòng lặp sẽ tự động Retry
			}
		}(userID, addr)
	}

	wg.Wait()

	logger.Info("🚀 [UDP-RELIABLE-BROADCAST-RESULT]",
		"clients", len(clientsToBroadcast),
		"delivered_with_ack", successCount,
		"duration", time.Since(start).String(),
	)
}

func (p *benchmarkPool) ProcessAck(notificationID string, userID string) {
	ackKey := notificationID + ":" + userID
	p.acks.Store(ackKey, true)
}
