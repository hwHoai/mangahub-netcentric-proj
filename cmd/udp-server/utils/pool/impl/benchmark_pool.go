package udp_pool_impl

import (
	"crypto/sha256"
	"encoding/json"
	"mangahub/pkg/logger"
	"mangahub/pkg/types"
	"net"
	"sync"
	"sync/atomic"
	"time"

	udp_pools "mangahub/cmd/udp-server/utils/pool"
)

type benchmarkPool struct {
	clients map[string]*net.UDPAddr
	acks    sync.Map // Lưu trữ trạng thái ACK: key là "notifID:userID"
	mu      sync.RWMutex
}

var _ udp_pools.UDPPool = (*benchmarkPool)(nil)

func NewBenchmarkPool() udp_pools.UDPPool {
	return &benchmarkPool{
		clients: make(map[string]*net.UDPAddr),
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
	p.clients[userID] = addr
}

func (p *benchmarkPool) Unregister(userID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.clients, userID)
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
