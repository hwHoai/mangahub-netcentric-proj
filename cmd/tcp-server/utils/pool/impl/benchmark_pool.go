package pool_impl

import (
	benchmarks_prometheus "mangahub/benchmarks/prometheus"
	"mangahub/pkg/logger"
	"net"
	"sync"

	pool "mangahub/cmd/tcp-server/utils/pool"
)

type BenchmarkPool struct {
	mu      sync.RWMutex
	clients map[string][]net.Conn
	metrics *benchmarks_prometheus.Metrics
}

var _ pool.ConnectionPool = (*BenchmarkPool)(nil)

func NewBenchmarkPool(metrics *benchmarks_prometheus.Metrics) *BenchmarkPool {
	return &BenchmarkPool{
		clients: make(map[string][]net.Conn),
		metrics: metrics,
	}
}

func (p *BenchmarkPool) Register(userID string, conn net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clients[userID] = append(p.clients[userID], conn)
	p.metrics.ActiveConnections.Inc()
	p.metrics.TotalRequests.Inc()
	p.metrics.ResponsesSent.Inc()
	logger.Info("TCP Subscriber registered (MOCK)", "userID", userID, "total_conns", len(p.clients[userID]))
}

func (p *BenchmarkPool) Unregister(conn net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for userID, conns := range p.clients {
		for i, c := range conns {
			if c == conn {
				p.clients[userID] = append(conns[:i], conns[i+1:]...)
				logger.Info("TCP Subscriber session closed (MOCK)", "addr", conn.RemoteAddr().String(), "userID", userID)
				if len(p.clients[userID]) == 0 {
					delete(p.clients, userID)
				}

				p.metrics.ActiveConnections.Dec()
				return
			}
		}
	}
}

func (p *BenchmarkPool) Broadcast(userID string, message []byte) {
	p.mu.RLock()
	connsToBroadcast := make([]net.Conn, 0)
	if conns, exists := p.clients[userID]; exists {
		connsToBroadcast = make([]net.Conn, len(conns))
		copy(connsToBroadcast, conns)
	}
	p.mu.RUnlock()

	if len(connsToBroadcast) == 0 {
		return
	}

	for _, conn := range connsToBroadcast {
		_, err := conn.Write(append(message, '\n'))
		if err != nil {
			logger.Error("Error sending to benchmark connection", "error", err)
		}
		p.metrics.ResponsesSent.Inc()
	}

	p.metrics.TotalRequests.Inc()
}

func (p *BenchmarkPool) BroadcastOne(userID string, message []byte, conn net.Conn) {
	// For benchmark, we just write to the connection
	conn.Write(append(message, '\n'))
}

func (p *BenchmarkPool) countTotal() int {
	total := 0
	for _, conns := range p.clients {
		total += len(conns)
	}
	return total
}
