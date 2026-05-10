package pool_impl

import (
	"context"
	benchmarks_prometheus "mangahub/benchmarks/prometheus"
	"mangahub/pkg/logger"
	"net"
	"sync"
	"time"

	pool "mangahub/cmd/tcp-server/utils/pool"
	"mangahub/proto/user_manga"
)

type ChapterSyncPool struct {
	mu           sync.RWMutex
	clients      map[string][]net.Conn
	lastChapters map[string]string
	grpcClient   user_manga.GRPCUserMangaServiceClient
	metrics      *benchmarks_prometheus.Metrics
}

var _ pool.ConnectionPool = (*ChapterSyncPool)(nil)

func NewChapterSyncPool(grpcClient user_manga.GRPCUserMangaServiceClient, metrics *benchmarks_prometheus.Metrics) *ChapterSyncPool {
	return &ChapterSyncPool{
		clients:      make(map[string][]net.Conn),
		lastChapters: make(map[string]string),
		grpcClient:   grpcClient,
		metrics:      metrics,
	}
}

func (p *ChapterSyncPool) Register(userID string, conn net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Avoid duplicate connections for the same user
	for _, c := range p.clients[userID] {
		if c == conn {
			return
		}
	}

	p.clients[userID] = append(p.clients[userID], conn)
	p.metrics.ActiveConnections.Inc()
	p.metrics.TotalRequests.Inc()
	p.metrics.ResponsesSent.Inc()
}

func (p *ChapterSyncPool) Unregister(conn net.Conn) {
	p.mu.Lock()
	var disconnectedUserID string
	for userID, conns := range p.clients {
		for i, c := range conns {
			if c == conn {
				p.clients[userID] = append(conns[:i], conns[i+1:]...)
				if len(p.clients[userID]) == 0 {
					disconnectedUserID = userID
					delete(p.clients, userID)
				}
				p.metrics.ActiveConnections.Dec()
				goto done
			}
		}
	}
done:
	p.mu.Unlock()

	if disconnectedUserID != "" {
		p.mu.RLock()
		chapterID, ok := p.lastChapters[disconnectedUserID]
		p.mu.RUnlock()

		if ok {
			go p.saveReadingProgress(disconnectedUserID, chapterID)
		}
	}
}

func (p *ChapterSyncPool) saveReadingProgress(userID, chapterID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := p.grpcClient.StoreReadingProgress(ctx, &user_manga.StoreReadingProgressRequest{
		UserId:    userID,
		ChapterId: chapterID,
	})
	if err != nil {
		logger.Error("Failed to save reading progress", "userID", userID, "chapterID", chapterID, "error", err)
		return
	}
	logger.Info("Successfully saved reading progress", "userID", userID, "chapterID", chapterID)

	p.mu.Lock()
	delete(p.lastChapters, userID)
	p.mu.Unlock()
}

func (p *ChapterSyncPool) Broadcast(userID string, message []byte) {
	p.mu.RLock()
	conns, exists := p.clients[userID]
	p.mu.RUnlock()

	if !exists {
		return
	}

	for _, conn := range conns {
		_, err := conn.Write(append(message, '\n'))
		if err != nil {
			logger.Error("Error sending to connection", "error", err)
		}
		p.metrics.ResponsesSent.Inc()
	}

	p.metrics.TotalRequests.Inc()
}

func (p *ChapterSyncPool) BroadcastOne(userID string, message []byte, conn net.Conn) {
	p.mu.RLock()
	_, exists := p.clients[userID]
	p.mu.RUnlock()

	if !exists {
		return
	}

	conn.Write(append(message, '\n'))
	p.metrics.ResponsesSent.Inc()
	p.metrics.TotalRequests.Inc()
}

func (p *ChapterSyncPool) UpdateLastChapter(userID string, chapterID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.lastChapters[userID] = chapterID

	p.metrics.TotalRequests.Inc()
	p.metrics.ResponsesSent.Inc()
}

func (p *ChapterSyncPool) GetLastChapter(userID string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	p.metrics.TotalRequests.Inc()
	p.metrics.ResponsesSent.Inc()
	return p.lastChapters[userID]
}
