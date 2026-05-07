package handler

import (
	"encoding/json"
	"net"

	"mangahub/cmd/udp-server/dispatch"
	"mangahub/cmd/udp-server/utils/pool"
	"mangahub/pkg/logger"
	"mangahub/pkg/types"
)

type NotificationHandler struct {
	chapterPool udp_pools.UDPPool
	messagePool udp_pools.UDPPool
}

func NewNotificationHandler(chapterPool, messagePool udp_pools.UDPPool) *NotificationHandler {
	return &NotificationHandler{
		chapterPool: chapterPool,
		messagePool: messagePool,
	}
}

func (h *NotificationHandler) ClientRegisterHandler(s *dispatch.UDPServer, clientAddr *net.UDPAddr, payload types.UDPMessage) {
	logger.Info("Received client register request", "addr", clientAddr.String())
	var data struct {
		UserID string `json:"user_id"`
	}

	if err := json.Unmarshal(payload.Payload, &data); err != nil {
		logger.Error("Error unmarshaling client_register payload", "error", err)
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
		logger.Error("Error unmarshaling broadcast_chapter payload", "error", err)
		return
	}

	logger.Info("Broadcasting new chapter", "chapterID", data.ChapterID, "mangaID", data.MangaID)

	message := map[string]interface{}{
		"action":  "chapter:on_new_chapter_notification",
		"payload": data,
	}

	h.chapterPool.Broadcast(s.Conn, data.MangaID, message)
}

func (h *NotificationHandler) BroadcastMessageHandler(s *dispatch.UDPServer, clientAddr *net.UDPAddr, payload types.UDPMessage) {
	var data types.NewMessageNotificationPayload

	if err := json.Unmarshal(payload.Payload, &data); err != nil {
		logger.Error("Error unmarshaling broadcast_message payload", "error", err)
		return
	}

	logger.Info("Broadcasting new message notification", "roomID", data.RoomID, "sender", data.SenderName)

	message := map[string]interface{}{
		"action":  "chat:on_new_message_notification",
		"payload": data,
	}

	h.messagePool.Broadcast(s.Conn, data.RoomID, message)
}

func (h *NotificationHandler) NotificationAckHandler(s *dispatch.UDPServer, clientAddr *net.UDPAddr, payload types.UDPMessage) {
	var data types.NotificationAckPayload

	if err := json.Unmarshal(payload.Payload, &data); err != nil {
		logger.Error("Error unmarshaling notification_ack payload", "error", err)
		return
	}

	logger.Info("Received notification ACK", "notificationID", data.NotificationID, "userID", data.UserID)

	// Try processing ACK in both pools (one will find it)
	h.chapterPool.ProcessAck(data.NotificationID, data.UserID)
	h.messagePool.ProcessAck(data.NotificationID, data.UserID)
}
