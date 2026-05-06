package websocket_impl

import (
	"context"
	"log"
	"time"

	ws_utils_pool "mangahub/cmd/websocket-server/utils/pool"
	udp_services "mangahub/internal/udp"
	ws_service "mangahub/internal/websocket"
	"mangahub/proto/manga"
	"mangahub/proto/message"

	"github.com/gorilla/websocket"
)

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
		log.Printf("Failed to validate room %s (exists: %v): %v", roomID, checkRes.GetExists(), err)
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

	for {
		var incoming struct {
			Content string `json:"content"`
		}
		if err := conn.ReadJSON(&incoming); err != nil {
			log.Printf("Error reading message from %s: %v", userID, err)
			break
		}
		log.Printf("Received message from %s: %s", userID, incoming.Content)

		// 1. Broadcast to pool
		log.Printf("Broadcasting message to room %s", roomID)
		msg := ws_utils_pool.Message{
			RoomID:  roomID,
			Content: incoming.Content,
			Sender:  userID,
		}
		s.pool.Broadcast(msg)
		log.Printf("Broadcast completed")

		go func() {
			// 2. Save to gRPC
			log.Printf("Saving message to gRPC for room %s", roomID)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err := s.messageClient.SaveMessage(ctx, &message.SaveMessageRequest{
				SenderId: userID,
				RoomId:   roomID,
				Content:  incoming.Content,
			})
			cancel()

			if err != nil {
				log.Printf("Failed to save message: %v", err)
				// Continue even if save fails? User choice.
			} else {
				log.Printf("Message saved successfully")
				// 3. Send UDP Notification to users reading in this room
			if s.udpClient != nil {
				log.Printf("Sending UDP notification for room %s", roomID)
				err := s.udpClient.SendNewMessageNotification(roomID, userID, incoming.Content)
				if err != nil {
					log.Printf("UDP Notification failed: %v", err)
				} else {
					log.Printf("UDP Notification sent")
				}
			}
			}
		}()
	}
}
