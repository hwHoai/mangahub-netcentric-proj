package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"mangahub/cmd/udp-server/dispatch"
	"mangahub/pkg/types"
	udp_pool "mangahub/cmd/udp-server/utils/pools"
)

type NotificationHandler struct {
	chapterPool udp_pool.UDPPool
	messagePool udp_pool.UDPPool
}

func NewNotificationHandler(chapterPool, messagePool udp_pool.UDPPool) *NotificationHandler {
	return &NotificationHandler{
		chapterPool: chapterPool,
		messagePool: messagePool,
	}
}

func (h *NotificationHandler) ClientRegisterHandler(s *dispatch.UDPServer, clientAddr *net.UDPAddr, payload types.UDPMessage) {
	fmt.Printf("Received client register request from %s: %+v\n", clientAddr, payload)
	var data struct {
		UserID string `json:"user_id"`
	}

	if err := json.Unmarshal(payload.Payload, &data); err != nil {
		fmt.Printf("Error unmarshaling client_register payload: %v\n", err)
		return
	}

	if data.UserID != "" {
		h.chapterPool.Register(data.UserID, clientAddr)
		h.messagePool.Register(data.UserID, clientAddr)
	}
}

func (h *NotificationHandler) BroadcastChapterHandler(s *dispatch.UDPServer, clientAddr *net.UDPAddr, payload types.UDPMessage) {
	var data types.NewChapterNotificationPayload

	if err := json.Unmarshal(payload.Payload, &data); err != nil {
		fmt.Printf("Error unmarshaling broadcast_chapter payload: %v\n", err)
		return
	}

	log.Printf("Broadcasting new chapter %s for manga %s", data.ChapterID, data.MangaID)

	message := map[string]interface{}{
		"action": "chapter:on_new_chapter_notification",
		"payload": data,
	}

	h.chapterPool.Broadcast(s.Conn, data.MangaID, message)
}

func (h *NotificationHandler) BroadcastMessageHandler(s *dispatch.UDPServer, clientAddr *net.UDPAddr, payload types.UDPMessage) {
	var data types.NewMessageNotificationPayload

	if err := json.Unmarshal(payload.Payload, &data); err != nil {
		fmt.Printf("Error unmarshaling broadcast_message payload: %v\n", err)
		return
	}

	log.Printf("Broadcasting new message in room %s from %s", data.RoomID, data.SenderName)

	message := map[string]interface{}{
		"action": "chat:on_new_message_notification",
		"payload": data,
	}

	h.messagePool.Broadcast(s.Conn, data.RoomID, message)
}
