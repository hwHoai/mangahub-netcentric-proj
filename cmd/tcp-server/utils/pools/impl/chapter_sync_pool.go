package pool_impl

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pools "mangahub/cmd/tcp-server/utils/pools"
	"mangahub/proto/user_manga"
)

type ChapterSyncPool struct {
	mu           sync.RWMutex
	clients      map[string][]net.Conn
	lastChapters map[string]string
	grpcClient   user_manga.GRPCUserMangaServiceClient
}

var _ pools.ConnectionPool = (*ChapterSyncPool)(nil)

func NewChapterSyncPool(grpcClient user_manga.GRPCUserMangaServiceClient) *ChapterSyncPool {
	return &ChapterSyncPool{
		clients:      make(map[string][]net.Conn),
		lastChapters: make(map[string]string),
		grpcClient:   grpcClient,
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
		log.Printf("Failed to save reading progress for user %s, chapter %s: %v", userID, chapterID, err)
		return
	}
	log.Printf("Successfully saved reading progress for user %s, chapter %s", userID, chapterID)

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
			fmt.Printf("Error sending to connection: %v\n", err)
		}
	}
}

func (p *ChapterSyncPool) BroadcastOne(userID string, message []byte, conn net.Conn) {
	p.mu.RLock()
	_, exists := p.clients[userID]
	p.mu.RUnlock()

	if !exists {
		return
	}

	conn.Write(append(message, '\n'))
}

func (p *ChapterSyncPool) UpdateLastChapter(userID string, chapterID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.lastChapters[userID] = chapterID
}

func (p *ChapterSyncPool) GetLastChapter(userID string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.lastChapters[userID]
}
