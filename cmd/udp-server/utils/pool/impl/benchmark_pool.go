package udp_pool_impl

import (
	"encoding/json"
	"mangahub/pkg/logger"
	"maps"
	"net"
	"sync"
	"time"

	udp_pools "mangahub/cmd/udp-server/utils/pool"
)

type benchmarkPool struct {
	mu      sync.RWMutex
	clients map[string]*net.UDPAddr
}

var _ udp_pools.UDPPool = (*benchmarkPool)(nil)

func NewBenchmarkPool() udp_pools.UDPPool {
	return &benchmarkPool{
		clients: make(map[string]*net.UDPAddr),
	}
}

func (p *benchmarkPool) Register(userID string, addr *net.UDPAddr) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clients[userID] = addr
	logger.Info("UDP benchmark client registered", "userID", userID, "addr", addr.String(), "total", len(p.clients))
}

func (p *benchmarkPool) Unregister(userID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.clients, userID)
	logger.Info("UDP benchmark client unregistered", "userID", userID, "remaining", len(p.clients))
}

func (p *benchmarkPool) Broadcast(conn *net.UDPConn, id string, payload map[string]interface{}) {
	p.mu.RLock()
	clientsToBroadcast := make(map[string]*net.UDPAddr, len(p.clients))
	maps.Copy(clientsToBroadcast, p.clients)
	p.mu.RUnlock()

	dataBytes, _ := json.Marshal(payload)
	dataBytes = append(dataBytes, '\n')

	start := time.Now()
	for _, addr := range clientsToBroadcast {
		conn.WriteToUDP(dataBytes, addr)
	}
	
	logger.Info("Dispatched MOCK Chapter Notification to subscribers", 
		"clients", len(clientsToBroadcast), 
		"time", time.Since(start).String(),
	)
}

func (p *benchmarkPool) ProcessAck(notificationID string, userID string) {
	// Not used for benchmark
}
