package websocket_impl

import (
	"context"
	"sync"
	"time"

	ws_utils_pool "mangahub/cmd/websocket-server/utils/pool"
	udp_services "mangahub/internal/udp"
	ws_service "mangahub/internal/websocket"
	"mangahub/pkg/logger"
	"mangahub/proto/manga"
	"mangahub/proto/message"

	"github.com/gorilla/websocket"
)

// rateLimiter implements a simple token bucket rate limiter per connection.
type rateLimiter struct {
	mu         sync.Mutex
	tokens     int
	maxTokens  int
	refillRate time.Duration
	lastRefill time.Time
}

func newRateLimiter(maxTokens int, refillRate time.Duration) *rateLimiter {
	return &rateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// allow checks if a message is allowed under the rate limit.
func (rl *rateLimiter) allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)
	tokensToAdd := int(elapsed / rl.refillRate)
	if tokensToAdd > 0 {
		rl.tokens += tokensToAdd
		if rl.tokens > rl.maxTokens {
			rl.tokens = rl.maxTokens
		}
		rl.lastRefill = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

type WSChatServiceImpl struct {
	pool          ws_utils_pool.ChatPool
	messageClient message.GRPCMessageServiceClient
	mangaClient   manga.GRPCMangaServiceClient
	udpClient     udp_services.UDPChapterNotificationServices
}

var _ ws_service.ChatService = (*WSChatServiceImpl)(nil)

func NewWSChatService(
	pool ws_utils_pool.ChatPool,
	messageClient message.GRPCMessageServiceClient,
	mangaClient manga.GRPCMangaServiceClient,
	udpClient udp_services.UDPChapterNotificationServices,
) ws_service.ChatService {
	return &WSChatServiceImpl{
		pool:          pool,
		messageClient: messageClient,
		mangaClient:   mangaClient,
		udpClient:     udpClient,
	}
}

func (s *WSChatServiceImpl) HandleWSChatTunnel(conn *websocket.Conn, userID string, roomID string) {
	// 0. Validate Room ID (Manga Existence)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checkRes, err := s.mangaClient.CheckMangaExists(ctx, &manga.CheckMangaExistsRequest{Id: roomID})
	if err != nil || !checkRes.Exists {
		logger.Error("Failed to validate room", "roomID", roomID, "exists", checkRes.GetExists(), "error", err)
		conn.WriteJSON(map[string]interface{}{"error": "invalid room (manga not found)"})
		conn.Close()
		return
	}

	client := &ws_utils_pool.Client{
		UserID: userID,
		RoomID: roomID,
		Conn:   conn,
		Pool:   s.pool,
	}

	s.pool.Register(client)

	defer func() {
		s.pool.Unregister(client)
		conn.Close()
	}()

	// Rate limiter: 10 messages per 500ms refill rate (effectively 20 msg/s burst, 2 msg/s sustained)
	limiter := newRateLimiter(10, 500*time.Millisecond)

	for {
		var incoming struct {
			Content string `json:"content"`
		}
		if err := conn.ReadJSON(&incoming); err != nil {
			logger.Error("Error reading message", "userID", userID, "error", err)
			break
		}

		// Check rate limit before processing
		if !limiter.allow() {
			conn.WriteJSON(map[string]interface{}{
				"error":  "rate_limit_exceeded",
				"sender": "system",
			})
			continue
		}

		logger.Info("Received message", "userID", userID, "content", incoming.Content)

		// 1. Broadcast to pool
		logger.Info("Broadcasting message", "roomID", roomID)
		msg := ws_utils_pool.Message{
			RoomID:  roomID,
			Content: incoming.Content,
			Sender:  userID,
		}
		s.pool.Broadcast(msg)
		logger.Info("Broadcast completed")

		go func() {
			// 2. Save to gRPC
			logger.Info("Saving message to gRPC", "roomID", roomID)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := s.messageClient.SaveMessage(ctx, &message.SaveMessageRequest{
				SenderId: userID,
				RoomId:   roomID,
				Content:  incoming.Content,
			})
			cancel()

			if err != nil {
				logger.Error("Failed to save message", "error", err)
			} else {
				logger.Info("Message saved successfully")
				// 3. Send UDP Notification to users reading in this room
			if s.udpClient != nil {
				logger.Info("Sending UDP notification", "roomID", roomID)
				err := s.udpClient.SendNewMessageNotification(roomID, userID, incoming.Content)
				if err != nil {
					logger.Error("UDP Notification failed", "error", err)
				} else {
					logger.Info("UDP Notification sent")
				}
			}
			}
		}()
	}
}

