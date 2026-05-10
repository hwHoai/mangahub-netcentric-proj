package handler

import (
	"encoding/json"
	pool_impl "mangahub/cmd/tcp-server/utils/pool/impl"
	"mangahub/pkg/logger"
	"net"
)

type ChapterSyncHandler struct {
	pool *pool_impl.ChapterSyncPool
}

func NewChapterSyncHandler(pool *pool_impl.ChapterSyncPool) *ChapterSyncHandler {
	return &ChapterSyncHandler{pool: pool}
}

func (h *ChapterSyncHandler) RegisterConnectionHandler(conn net.Conn, payload any) {
	var data struct {
		UserID string `json:"user_id"`
	}
	
	raw, ok := payload.(json.RawMessage)
	if !ok {
		logger.Error("Invalid payload type for register")
		return
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		logger.Error("Error unmarshaling register payload", "error", err)
		return
	}

	h.pool.Register(data.UserID, conn)

	newProgressMsg, _ := json.Marshal(map[string]any{
		"action": "chapter_sync:on_sync_progress",
		"payload": map[string]string{
			"chapter_id": h.pool.GetLastChapter(data.UserID),
		},
	})

	h.pool.BroadcastOne(data.UserID, newProgressMsg, conn)
	logger.Info("User registered connection", "userID", data.UserID)
	conn.Write([]byte("Registered successfully\n"))
}

func (h *ChapterSyncHandler) BroadcastReadHandler(conn net.Conn, payload any) {
	var data struct {
		UserID    string `json:"user_id"`
		ChapterID string `json:"chapter_id"`
	}

	raw, ok := payload.(json.RawMessage)
	if !ok {
		logger.Error("Invalid payload type for broadcast_read")
		return
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		logger.Error("Error unmarshaling broadcast_read payload", "error", err)
		return
	}

	logger.Info("Broadcasting sync_reading", "userID", data.UserID, "chapterID", data.ChapterID)
	
	h.pool.UpdateLastChapter(data.UserID, data.ChapterID)

	syncMsg, _ := json.Marshal(map[string]any{
		"action": "chapter_sync:on_sync_progress",
		"payload": map[string]string{
			"chapter_id": data.ChapterID,
		},
	})
	
	h.pool.Broadcast(data.UserID, syncMsg)
}
